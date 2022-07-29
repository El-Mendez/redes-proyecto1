package query

import "encoding/xml"

type UnregisterQuery struct {
	XMLName xml.Name `xml:"jabber:iq:register query"`
	Remove  string   `xml:"remove"`
}

func (q *UnregisterQuery) isQuery() {}

// As this response does not get sent by the server, there's no need to add a namespace.
