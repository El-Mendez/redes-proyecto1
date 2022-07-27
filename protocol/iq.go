package protocol

import (
	"encoding/xml"
	"math/rand"
)

type IQ struct {
	XMLName  xml.Name `xml:"iq"`
	ID       string   `xml:"id,attr"`
	Type     string   `xml:"type,attr"`
	To       string   `xml:"to,attr,omitempty"`
	From     string   `xml:"from,attr,omitempty"`
	Contents string   `xml:",innerxml"`
	//Error    *Error   `xml:"error"`
}

func (iq *IQ) addContents(v interface{}) error {
	data, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	iq.Contents = string(data)
	return nil
}

func (iq *IQ) getContents(v interface{}) error {
	return xml.Unmarshal([]byte(iq.Contents), v)
}

// GenerateID - Recovered from https://golangdocs.com/generate-random-string-in-golang
func GenerateID() string {
	CHARACTERS := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	s := make([]rune, 10)
	for i := range s {
		s[i] = CHARACTERS[rand.Intn(len(CHARACTERS))]
	}
	return string(s)
}

// ---- the specific iq body types

type bindQuery struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	JID      string   `xml:"jid,omitempty"`
}

type rosterQuery struct {
	XMLName     xml.Name     `xml:"jabber:iq:roster query"`
	RosterItems []rosterItem `xml:"item"`
}
type rosterItem struct {
	Jid string `xml:"jid,attr"`
}
