package query

import "encoding/xml"

var CLOSE_IBB_XML_NAME = xml.Name{"http://jabber.org/protocol/ibb", "close"}

type CloseIBBQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/ibb close"`
	Sid     string   `xml:"sid,attr"`
}

func (q *CloseIBBQuery) isQuery() {}
