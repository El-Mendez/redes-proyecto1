package protocol

import "encoding/xml"

type features struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams features"`
}
