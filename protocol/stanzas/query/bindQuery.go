package query

import "encoding/xml"

var BIND_QUERY_XML_NAME = xml.Name{"urn:ietf:params:xml:ns:xmpp-bind", "bind"}

type BindQuery struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	JID      string   `xml:"jid,omitempty"`
}

func (q *BindQuery) isQuery() {}
