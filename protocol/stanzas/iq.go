package stanzas

import (
	"encoding/xml"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	"math/rand"
)

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	ID      string   `xml:"id,attr"`
	Type    string   `xml:"type,attr"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	Query   query.Query
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

func (iq *IQ) isStanza() {}

func (iq *IQ) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Get the IQ attributes
	iq.XMLName = start.Name
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			iq.ID = attr.Value
		case "type":
			iq.Type = attr.Value
		case "to":
			iq.To = attr.Value
		case "from":
			iq.From = attr.Value
		}
	}

	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		var q query.Query

		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name {
			case query.BIND_QUERY_XML_NAME:
				q = &query.BindQuery{}
			case query.REGISTER_QUERY_XML_NAME:
				q = &query.RegisterQuery{}
			case query.ROSTER_QUERY_XML_NAME:
				q = &query.RosterQuery{}
			}
			// known child element found, decode it
			if q != nil {
				if err := d.DecodeElement(q, &tt); err != nil {
					return err
				}
				iq.Query = q
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
	return nil
}
