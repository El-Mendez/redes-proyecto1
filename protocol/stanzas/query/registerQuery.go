package query

import "encoding/xml"

// REGISTER_QUERY_XML_NAME represents the xml namespace corresponding to XEP-0077. Used when unmarshalling an IQ.
var REGISTER_QUERY_XML_NAME = xml.Name{"jabber:iq:register", "query"}

// RegisterQuery represents a Data Transfer Object when Unmarshalling a Stanza.
type RegisterQuery struct {
	XMLName  xml.Name `xml:"jabber:iq:register query"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Email    string   `xml:"email"`
}

func (q *RegisterQuery) isQuery() {}
