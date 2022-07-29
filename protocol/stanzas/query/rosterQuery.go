package query

import "encoding/xml"

var ROSTER_QUERY_XML_NAME = xml.Name{"jabber:iq:roster", "query"}

type RosterQuery struct {
	XMLName     xml.Name     `xml:"jabber:iq:roster query"`
	RosterItems []RosterItem `xml:"item"`
}
type RosterItem struct {
	Jid string `xml:"jid,attr"`
}

func (q *RosterQuery) isQuery() {}
