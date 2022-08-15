package query

import "encoding/xml"

// ROSTER_QUERY_XML_NAME represents the xml namespace corresponding to roster queries. Only used when Unmarshalling an Stanza.
var ROSTER_QUERY_XML_NAME = xml.Name{"jabber:iq:roster", "query"}

// RosterQuery is a DTO for request regarding the roster.
type RosterQuery struct {
	XMLName     xml.Name     `xml:"jabber:iq:roster query"`
	RosterItems []RosterItem `xml:"item"`
}

// RosterItem is just a wrapper for the user jid when unmarshalling a user. Created to bypass some limitations of the
// encoding/xml library to unmarshal a xml tag attribute without the need to rewrite an Unmarshalling function again.
type RosterItem struct {
	Jid string `xml:"jid,attr"`
}

func (q *RosterQuery) isQuery() {}
