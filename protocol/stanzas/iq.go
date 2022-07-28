package stanzas

import (
	"encoding/xml"
	"github.com/el-mendez/redes-proyecto1/protocol/stanzas/query"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"math/rand"
)

type IQ struct {
	XMLName  xml.Name `xml:"iq"`
	ID       string   `xml:"id,attr"`
	Type     string   `xml:"type,attr"`
	To       string   `xml:"to,attr,omitempty"`
	From     string   `xml:"from,attr,omitempty"`
	Contents string   `xml:",innerxml"`
	//Error    *Error   `xml:"error"`
}
type IQParams struct {
	ID    string
	Type  string
	To    string
	From  string
	Query query.Query
}

func NewIQ(params *IQParams) *IQ {
	iq := IQ{
		ID:   params.ID,
		Type: params.Type,
		To:   params.To,
		From: params.From,
	}
	utils.Successful(iq.AddContents(params.Query), "Could not embed IQ with Query: %v")
	return &iq
}

func (iq *IQ) AddContents(q query.Query) error {
	data, err := xml.Marshal(q)
	if err != nil {
		return err
	}
	iq.Contents = string(data)
	return nil
}

func (iq *IQ) GetContents(q query.Query) error {
	return xml.Unmarshal([]byte(iq.Contents), q)
}

// GenerateID - Recovered from https://golangdocs.com/generate-random-string-in-golang
func GenerateID() string {
	CHARACTERS := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	s := make([]rune, 10)
	for i := range s {
		s[i] = CHARACTERS[rand.Intn(len(CHARACTERS))]
	}
	return string(s)
}

func (iq *IQ) isStanza() bool {
	return true
}
