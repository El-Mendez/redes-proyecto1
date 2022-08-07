package protocol

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"net"
)

type Client struct {
	stream  *Stream
	jid     *JID
	Send    chan<- stanzas.Stanza
	Receive <-chan stanzas.Stanza

	outgoing <-chan stanzas.Stanza
	incoming chan<- stanzas.Stanza
	isClosed utils.AtomicBool
}

func (client *Client) FullJid() string {
	return client.jid.String()
}

func (client *Client) BaseJid() string {
	return client.jid.BaseJid()
}

func SignUp(jid *JID, password string) (*Client, error) {
	utils.Logger.Info("Attempting to create channel and signup.")

	stream := MakeStream(jid.Domain)
	if stream == nil {
		return nil, fmt.Errorf("could not connect to %v", jid.Domain)
	}

	client := &Client{jid: jid, stream: stream}
	id := stanzas.GenerateID()

	var request stanzas.Stanza = &stanzas.IQ{
		ID:   id,
		Type: "set",
		Query: &query.RegisterQuery{
			Username: jid.Username,
			Password: password,
		},
	}
	utils.Successful(client.sendStanza(request), "Could not send signup Stanza.")

	response, err := client.getStanza()
	if err != nil {
		return nil, err
	}

	iq, ok := response.(*stanzas.IQ)
	if !ok || iq.ID != id {
		utils.Logger.Fatalf("Expected a login IQ stanza, instead got: %T=%v", response, response)
	}

	if iq.Type != "result" {
		utils.Logger.Infof("Could not create account: %v", jid.BaseJid())
		client.stream.Close()
		return nil, fmt.Errorf("username already exists")
	}

	return logIn(jid, password, stream)
}

func LogIn(jid *JID, password string) (*Client, error) {
	utils.Logger.Info("Attempting to create channel and login.")

	stream := MakeStream(jid.Domain)
	if stream == nil {
		return nil, fmt.Errorf("could not connect to %v", jid.Domain)
	}

	return logIn(jid, password, stream)
}

func logIn(jid *JID, password string, stream *Stream) (*Client, error) {
	toServer := make(chan stanzas.Stanza)
	fromServer := make(chan stanzas.Stanza)

	client := &Client{
		jid:     jid,
		stream:  stream,
		Send:    toServer,
		Receive: fromServer,

		outgoing: toServer,
		incoming: fromServer,
	}

	if err := client.authorize(password); err != nil {
		client.Close()
		return nil, err
	}

	utils.Successful(client.stream.Restart(jid.Domain), "could not restart stream after authorization successful: %v")

	client.bind()

	go client.pipeReceiving()
	go client.handleSending()

	return client, nil
}

func (client *Client) authorize(password string) error {
	secret := base64.StdEncoding.EncodeToString([]byte("\x00" + client.jid.Username + "\x00" + password))
	// this should be safe because base64 does not contain XML private characters
	err := client.stream.Write([]byte("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>" + secret + "</auth>"))
	if err != nil {
		utils.Logger.Errorf("Could not send login request")
		return err
	}
	utils.Logger.Debug("Sent LoginScreen request")

	// Check for response
	tag, _, err := client.stream.NextElement()
	if err != nil {
		return err
	}

	switch tag.Name.Local {
	case "success":
		utils.Logger.Info("Successfully logged in.")
		return nil
	case "failure":
		utils.Logger.Info("Could not log in. Incorrect password or account")
		return fmt.Errorf("not authorized")
	default:
		// idk, packages were lost I suppose
		utils.Logger.Errorf("Expected success/failure tag at log in but got: . %s", tag.Name)
		return fmt.Errorf("internet connection")
	}

}

func (client *Client) Close() {
	client.isClosed.Set(true)

	close(client.Send)
	close(client.incoming)

	client.stream.Close()
}

func (client *Client) SendMessage(to string, body string) {
	var message stanzas.Stanza = &stanzas.Message{
		Type: "chat",
		To:   to,
		From: client.jid.String(),
		Body: body,
	}

	client.Send <- message
}

func (client *Client) bind() {
	// Build the request IQ
	utils.Logger.Info("Attempting to bind")

	var request stanzas.Stanza = &stanzas.IQ{
		ID:   stanzas.GenerateID(),
		Type: "set",
		Query: &query.BindQuery{
			Resource: client.jid.DeviceName,
		},
	}
	utils.Successful(client.sendStanza(request), "Could not send bind request stanza.")

	// As binding only happens before logging in, we know the next Stanza must be the Binding response
	r, err := client.getStanza()
	if err != nil {
		utils.Logger.Fatalf("Could not bind the session: %v", err)
	}
	response, ok := r.(*stanzas.IQ)
	if !ok {
		utils.Logger.Fatalf("Expected a IQ as a binding response, instead got: %T", response)
	}

	bind, ok := response.Query.(*query.BindQuery)
	if !ok {
		utils.Logger.Fatalf("Expected the binded IQ response to contain <bind>, instead got: %T", response)
	}

	jid, ok := JIDFromString(bind.JID)
	if !ok {
		utils.Logger.Fatalf("Could not parse the server binded JID %v", bind.JID)
	}

	utils.Logger.Infof("Successfully binded as %v", jid.String())
	*client.jid = jid
}

func (client *Client) askRoster() {

}

func (client *Client) DeleteAccount() error {
	utils.Logger.Info("Attempting to delete account.")

	id := stanzas.GenerateID()

	var request stanzas.Stanza = &stanzas.IQ{
		ID:    id,
		Type:  "set",
		Query: &query.UnregisterQuery{},
	}
	client.Send <- request

	response, err := client.getStanza()
	if err != nil {
		return err
	}

	iq, ok := response.(*stanzas.IQ)
	if !ok || iq.ID != id {
		utils.Logger.Fatalf("Expected a response IQ stanza, instead got: %T=%v", response, response)
	}

	if iq.Type != "result" {
		utils.Logger.Infof("Could not delete account: %v", client.jid.BaseJid())
		client.Close()
		return fmt.Errorf("could not delete account")
	}

	utils.Logger.Infof("Succcessfully deleted %v account.", client.jid.BaseJid())
	return nil
}

func (client *Client) sendStanza(s stanzas.Stanza) error {
	data, err := xml.Marshal(s)
	if err != nil {
		utils.Logger.Fatal("Could not parse Stanza: %v", s)
	}
	return client.stream.Write(data)
}

func (client *Client) getStanza() (stanzas.Stanza, error) {
	tag, stanza, err := client.stream.NextElement()
	if err != nil {
		return nil, err
	}

	switch tag.Name.Local {
	case "iq":
		utils.Logger.Info("Received a IQ")
		iq := &stanzas.IQ{}
		utils.Successful(xml.Unmarshal(stanza, iq), "Could not unparse a query: %v")
		return iq, nil
	case "message":
		utils.Logger.Info("Received a message")
		message := &stanzas.Message{}
		utils.Successful(xml.Unmarshal(stanza, message), "Could not unparse message: %v")
		return message, nil
	case "presence":
		utils.Logger.Info("Received a presence")
		presence := &stanzas.Presence{}
		utils.Successful(xml.Unmarshal(stanza, presence), "Could not unparse presence: %v")
		return presence, nil
	default:
		utils.Logger.Errorf("Expected a iq/message tag, instead got: %v", tag.Name)
	}

	return nil, nil
}

func (client *Client) handleSending() {
	for s := range client.outgoing {
		err := client.sendStanza(s)
		if err != nil && !(errors.Is(err, net.ErrClosed) && client.isClosed.Get()) {
			utils.Logger.Errorf("Could not send %T stanza: %v", s, err)
		}
	}
}

func (client *Client) pipeReceiving() {
	for !client.isClosed.Get() {
		s, err := client.getStanza()
		if err != nil {
			if errors.Is(err, net.ErrClosed) && client.isClosed.Get() {
				utils.Logger.Info("Connection closed. Removing receiving pipe.")
				return
			} else {
				utils.Logger.Errorf("Could not receive stanza: %v", err)
			}
		}

		client.incoming <- s
	}
}
