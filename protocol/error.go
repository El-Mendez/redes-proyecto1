package protocol

import "encoding/xml"

type Error struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code,attr,omitempty"`
	Type    string   `xml:"type,attr"`
	Payload string   `xml:",innerxml"`
}
