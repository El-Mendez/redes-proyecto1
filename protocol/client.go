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

func SignIn(jid *JID, stream *Stream, password string) (*Client, error) {
	// Send the login request
	secret := base64.StdEncoding.EncodeToString([]byte("\x00" + jid.Username + "\x00" + password))
	err := stream.Write([]byte("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>" + secret + "</auth>"))
	if err != nil {
		return nil, err
	}
	utils.Logger.Debug("Sent Login request")

	// Check for response
	tag, _ := stream.NextElement()

	if tag.Name.Local == "success" {
		utils.Logger.Info("Successfully logged in.")
		return &Client{stream: stream, jid: jid}, nil

	} else if tag.Name.Local == "failure" {
		utils.Logger.Info("Could not log in.")
		return nil, fmt.Errorf("could not log in")
	}
	utils.Logger.Fatalf("Could not log in. %s", tag.Name)
	return nil, fmt.Errorf("found unexpected %s", tag.Name)
}
