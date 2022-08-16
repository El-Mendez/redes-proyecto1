package query

import "encoding/xml"

var OPEN_IBB_XML_NAME = xml.Name{"http://jabber.org/protocol/ibb", "open"}

type OpenIBBQuery struct {
	XMLName   xml.Name `xml:"http://jabber.org/protocol/ibb open"`
	BlockSize int      `xml:"block-size,attr"`
	Sid       string   `xml:"sid,attr"`
}

func (q *OpenIBBQuery) isQuery() {}
