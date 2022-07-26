package protocol

import (
	"bytes"
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
	_, e := stream.NextElement()
	return xml.Unmarshal(e, v)
}

func writeToken(enc *xml.Encoder, token xml.Token) {
	successful(enc.EncodeToken(token), "Could not encode token %v")
	successful(enc.Flush(), "Could not flush after writing token %v")
}

// NextElement This is just for logging purposes
func (stream *Stream) nextElement() (*xml.StartElement, []byte) {
	var temp struct {
		Inner []byte `xml:",innerxml"`
	}
	start, _ := stream.nextTag()
	end := start.End()
	successful(stream.decoder.DecodeElement(&temp, start), "Could not decode next xml element: %v")

	buffer := new(bytes.Buffer)
	enc := xml.NewEncoder(buffer)

	writeToken(enc, *start)
	buffer.Write(temp.Inner)
	writeToken(enc, end)

	return start, buffer.Bytes()
}

func (stream *Stream) NextElement() (*xml.StartElement, []byte) {
	tag, e := stream.nextElement()
	utils.Logger.Debugf("received: %v", string(e))

	return tag, e
}

func (stream *Stream) Skip() {
	_, e := stream.nextElement()
	utils.Logger.Debugf("skipped: %v", string(e))
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
	if err := stream.Write([]byte(xml.Header)); err != nil {
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

	tag, err := stream.nextTag()
	if err != nil && tag.Name != (xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"}) {
		return nil, fmt.Errorf("expected start tag")
	}

	utils.Logger.Info("Stream created successfully")

	feature := &features{}
	if err := stream.Read(feature); err != nil {
		utils.Logger.Errorf("Could not read features: %v", err)
	}

	return stream, nil
}

func (stream *Stream) Close() {
	utils.Logger.Info("Closing Stream")
	if err := stream.Write([]byte("</stream:stream>")); err != nil {
		utils.Logger.Warn("Could not close stream gracefully")
	}
}

func (stream *Stream) nextTag() (*xml.StartElement, error) {
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

func successful(err error, format string) {
	if err != nil {
		utils.Logger.Fatalf(format, err)
	}
}
