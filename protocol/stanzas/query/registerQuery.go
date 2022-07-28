package query

import "encoding/xml"

type RegisterQuery struct {
	XMLName  xml.Name `xml:"jabber:iq:register query"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
	Email    string   `xml:"email"`
}

func (q RegisterQuery) isQuery() bool {
	return true
}
