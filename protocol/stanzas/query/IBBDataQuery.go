package query

import "encoding/xml"

var DATA_IBB_XML_NAME = xml.Name{"http://jabber.org/protocol/ibb", "data"}

type IBBDataQuery struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/ibb data"`
	Value    string   `xml:",chardata"`
	Sequence uint16   `xml:"seq,attr"`
	Sid      string   `xml:"sid,attr"`
}

func (q *IBBDataQuery) isQuery() {}
