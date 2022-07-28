package protocol

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type Client struct {
	stream *Stream
	jid    *JID
}

func SignUp(jid *JID, password string) (*Client, error) {
	utils.Logger.Info("Attempting to create channel and signup.")

	stream := MakeStream(jid.Domain)
	if stream == nil {
		return nil, fmt.Errorf("could not connect to %v", jid.Domain)
	}

	client := &Client{jid: jid, stream: stream}
	id := stanzas.GenerateID()

	request := stanzas.NewIQ(&stanzas.IQParams{
		ID:   id,
		Type: "set",
		Query: &query.RegisterQuery{
			Username: jid.Username,
			Password: password,
		},
	})
	client.sendStanza(request)

	response := client.getStanza()
	iq, ok := response.(*stanzas.IQ)
	if !ok || iq.ID != id {
		utils.Logger.Fatalf("Expected a login IQ stanza, instead got: %T=%v", response)
	}

	if iq.Type != "result" {
		utils.Logger.Infof("Could not create account: %v", jid.BaseJid())
		client.Close()
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
	client := &Client{jid: jid, stream: stream}

	if err := client.authorize(password); err != nil {
		client.Close()
		return nil, err
	}

	utils.Successful(client.stream.Restart(jid.Domain), "could not restart stream after authorization successful: %v")

	client.bind()

	client.askRoster()
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
	utils.Logger.Debug("Sent Login request")

	// Check for response
	tag, _ := client.stream.NextElement()
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
	client.stream.Close()
}

func (client *Client) SendMessage(to string, body string) {
	message := stanzas.Message{
		Type: "chat",
		To:   to,
		From: client.jid.String(),
		Body: body,
	}

	client.sendStanza(&message)
}

func (client *Client) bind() {
	// Build the request IQ
	utils.Logger.Info("Attempting to bind")

	request := stanzas.NewIQ(&stanzas.IQParams{
		ID:   stanzas.GenerateID(),
		Type: "set",
		Query: &query.BindQuery{
			Resource: client.jid.DeviceName,
		},
	})
	client.sendStanza(request)

	response := client.getStanza()

	riq, ok := response.(*stanzas.IQ)
	if !ok {
		utils.Logger.Fatalf("Expected a IQ as a binding response, instead got: %T", response)
	}

	binded := query.BindQuery{}
	utils.Successful(riq.GetContents(&binded), "Expected the bind response Stanza to be IQBind: %v")

	new_jid, ok := JIDFromString(binded.JID)
	if !ok {
		utils.Logger.Fatalf("Could not parse the server binded JID %v", binded.JID)
	}

	utils.Logger.Infof("Successfully binded as %v", new_jid.String())
	*client.jid = new_jid

}

func (client *Client) askRoster() {
	// Build the request IQ
	utils.Logger.Info("Attempting to bind")

	request := stanzas.NewIQ(&stanzas.IQParams{
		ID:    stanzas.GenerateID(),
		Type:  "get",
		To:    client.jid.BaseJid(),
		From:  client.jid.String(),
		Query: &query.RosterQuery{},
	})
	client.sendStanza(request)

	r := client.getStanza()

	riq, ok := r.(*stanzas.IQ)
	if !ok {
		utils.Logger.Fatalf("Expected a IQ as a roster response, instead got: %T", r)
	}

	rbind := query.RosterQuery{}
	utils.Successful(riq.GetContents(&rbind), "Expected the roster response Stanza to be RosterQuery: %v")

	fmt.Println(rbind.RosterItems)
}

func (client *Client) sendStanza(s stanzas.Stanza) {
	data, err := xml.Marshal(s)
	if err != nil {
		utils.Logger.Fatal("Could not parse Stanza: %v", s)
	}
	utils.Successful(client.stream.Write(data), "Could not send stanza: %v")
}

func (client *Client) getStanza() stanzas.Stanza {
	tag, stanza := client.stream.NextElement()

	switch tag.Name.Local {
	case "iq":
		utils.Logger.Info("Received a IQ")
		iq := &stanzas.IQ{}
		utils.Successful(xml.Unmarshal(stanza, iq), "Could not unparse a query: %v")
		return iq
	case "message":
		utils.Logger.Info("Received a message")
		message := &stanzas.Message{}
		utils.Successful(xml.Unmarshal(stanza, message), "Could not unparse message: %v")
		return message
	default:
		utils.Logger.Errorf("Expected a iq/message tag, instead got: %v", tag.Name)
	}

	return nil
}
