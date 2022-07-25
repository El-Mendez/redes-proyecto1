package protocol

import (
	"encoding/xml"
	"fmt"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"io"
	"net"
)

const PORT = 5222

type Stream struct {
	conn    net.Conn
	decoder *xml.Decoder
}

func (stream *Stream) Write(data []byte) error {
	utils.Logger.Debugf("sending: %v", string(data))
	_, err := stream.conn.Write(data)
	return err
}

func (stream *Stream) Read(v interface{}) error {
	return stream.decoder.Decode(v)
}

func (stream *Stream) Skip() error {
	return stream.decoder.Skip()

}

// MakeStream creates a xmpp stream connected to a specific server.
func MakeStream(domain string) (*Stream, error) {
	address := fmt.Sprintf("%v:%v", domain, PORT)

	utils.Logger.Infof("Creating connection to %v", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		utils.Logger.Warnf("Could not create connection to %v", address)
		return nil, err
	}

	stream := &Stream{conn, xml.NewDecoder(conn)}

	// Start the server communication
	if err := stream.Write([]byte("<?xml version='1.0' encoding='utf-8'?>")); err != nil {
		utils.Logger.Warnf("Could send connection initiation to %v", domain)
		return nil, err
	}

	// Start the stream
	// TODO xml escape the domain address
	if err := stream.Write([]byte("<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' " +
		"to='" + domain + "' version='1.0'>")); err != nil {
		utils.Logger.Warnf("Could send stream invitation to %v", domain)
		return nil, err
	}

	tag, err := stream.NextTag()
	if err != nil && tag.Name != (xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"}) {
		return nil, fmt.Errorf("expected start tag")
	}

	utils.Logger.Info("Stream created successfully")

	feature := &features{}
	if err := stream.Read(feature); err != nil {
		utils.Logger.Errorf("Could not read features: %v", err)
	}

	utils.Logger.Info("Stream features ignored successfully")
	return stream, nil
}

func (stream *Stream) Close() {
	utils.Logger.Info("Closing Stream")
	if err := stream.Write([]byte("</stream:stream>")); err != nil {
		utils.Logger.Warn("Could not close stream gracefully")
	}
}

func (stream *Stream) NextTag() (*xml.StartElement, error) {
	for {
		token, err := stream.decoder.Token()
		if err != nil {
			return nil, err
		}

		switch tag := token.(type) {
		case xml.EndElement:
			return nil, io.EOF
		case xml.StartElement:
			return &tag, nil
		}
	}
}
