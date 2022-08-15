package stanzas

import "encoding/xml"

// Presence General use Presence Stanza.
type Presence struct {
	XMLName xml.Name `xml:"presence"`
	ID      string   `xml:"id,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	Show    string   `xml:"show,omitempty"`
	Status  []string `xml:"status,omitempty"`
}

func (p *Presence) isStanza() {}
