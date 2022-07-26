package protocol

import (
	"encoding/base64"
	"fmt"
	utils "github.com/el-mendez/redes-proyecto1/util"
)

type Client struct {
	stream *Stream
	jid    *JID
}

func SignIn(jid *JID, password string) (*Client, error) {
	stream := MakeStream(jid.Domain)
	if stream == nil {
		return nil, fmt.Errorf("could not connect to %v", jid.Domain)
	}

	client := &Client{jid: jid, stream: stream}

	if err := client.authorize(password); err != nil {
		client.Close()
		return nil, err
	}

	utils.Successful(client.stream.Restart(jid.Domain), "could not restart stream after authorization successful: %v")

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
