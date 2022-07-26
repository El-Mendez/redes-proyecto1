package stanzas

import "encoding/xml"

// Message basic message stanza structure.
type Message struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr,omitempty"`
	Type    string   `xml:"type,attr"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	Body    string   `xml:"body"`
}

func (message *Message) isStanza() {}
