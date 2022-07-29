package query

import "encoding/xml"

var REGISTER_QUERY_XML_NAME = xml.Name{"jabber:iq:register", "query"}

type RegisterQuery struct {
	XMLName  xml.Name `xml:"jabber:iq:register query"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Email    string   `xml:"email"`
}

func (q *RegisterQuery) isQuery() {}
