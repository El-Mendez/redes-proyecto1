package query

import "encoding/xml"

// BIND_QUERY_XML_NAME represents the xml namespace corresponding to the bind query. Used when unmarshalling an IQ.
var BIND_QUERY_XML_NAME = xml.Name{"urn:ietf:params:xml:ns:xmpp-bind", "bind"}

// BindQuery is just the dto inside an IQ for binding a session after logging in.
type BindQuery struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	JID      string   `xml:"jid,omitempty"`
}

func (q *BindQuery) isQuery() {}
