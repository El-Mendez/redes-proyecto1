package query

import "encoding/xml"

// UnregisterQuery is just a DTO for sending a request to delete the user account. As the server doesn't really forces
// the user to delete their account, we don't really need to unmarshal it. We only care about sending it.
type UnregisterQuery struct {
	XMLName xml.Name `xml:"jabber:iq:register query"`
	Remove  string   `xml:"remove"`
}

func (q *UnregisterQuery) isQuery() {}
