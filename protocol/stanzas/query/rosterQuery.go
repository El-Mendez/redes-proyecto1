package query

import "encoding/xml"

type RosterQuery struct {
	XMLName     xml.Name     `xml:"jabber:iq:roster query"`
	RosterItems []RosterItem `xml:"item"`
}
type RosterItem struct {
	Jid string `xml:"jid,attr"`
}

func (q RosterQuery) isQuery() bool {
	return true
}
