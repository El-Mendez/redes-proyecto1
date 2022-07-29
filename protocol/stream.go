package protocol

import (
	"bytes"
	"encoding/xml"
	"fmt"
	utils "github.com/el-mendez/redes-proyecto1/util"
	"io"
	"net"
	"time"
)

const PORT = 5222

type Stream struct {
	conn    net.Conn
	decoder *xml.Decoder
}

// Write sends a list of bytes through the stream to the server.
func (stream *Stream) Write(data []byte) error {
	utils.Logger.Debugf("sending: %v", string(data))
	_, err := stream.conn.Write(data)
	return err
}

// Read acts like xml.Unmarshal but for the next element in the stream.
func (stream *Stream) Read(v any) error {
	_, e := stream.NextElement()
	return xml.Unmarshal(e, v)
}

// NextElement returns the opening tag of the next XML element received and a []byte with the complete element.
func (stream *Stream) NextElement() (*xml.StartElement, []byte) {
	tag, e := stream.nextElement()
	utils.Logger.Debugf("received: %v", string(e))

	return tag, e
}

// MakeStream creates a xmpp stream connected to a specific server. Returns nil on initiation error.
func MakeStream(domain string) *Stream {
	address := fmt.Sprintf("%v:%v", domain, PORT)

	utils.Logger.Infof("Creating connection to %v", address)

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		utils.Logger.Warnf("Could not create connection to %v: %v", address, err)
		return nil
	}

	stream := &Stream{conn, xml.NewDecoder(conn)}

	// Start the server communication
	if err := stream.Write([]byte(xml.Header)); err != nil {
		utils.Logger.Warnf("Could not send xml header to server: %v", err)
		_ = conn.Close()
		return nil
	}

	// Start the stream with the xmpp server (the tags <stream:stream/>)
	if err := stream.Restart(domain); err != nil {
		utils.Logger.Errorf("Could not start stream at initiation: %v", err)
		_ = conn.Close()
		return nil
	}

	return stream
}

// Restart recreates a stream with the server when the server state resets.
func (stream *Stream) Restart(domain string) error {
	// Start the stream
	if err := stream.Write([]byte("<stream:stream xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' " +
		"to='" + domain + "' version='1.0'>")); err != nil {
		utils.Logger.Errorf("Could not send stream opening tag: %v", err)
		return err
	}

	tag, err := stream.nextTag()
	if err != nil && tag.Name != (xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"}) {
		utils.Logger.Errorf("Did not get stream opening tag: %v", err)
		return fmt.Errorf("expected start tag")
	}

	utils.Logger.Info("Stream restarted successfully")

	// I don't really know why, but without this it gets stuck in an endless loop. Might be a bug on my xml library.
	type features struct {
		XMLName xml.Name `xml:"http://etherx.jabber.org/streams features"`
	}

	feature := &features{}
	if err := stream.Read(feature); err != nil {
		utils.Logger.Errorf("Could not read features: %v", err)
	}

	return nil
}

// Close gracefully closes the XMPP connection.
func (stream *Stream) Close() {
	utils.Logger.Info("Closing Stream")
	if err := stream.Write([]byte("</stream:stream>")); err != nil {
		utils.Logger.Warn("Could not close stream gracefully")
	}
	utils.Successful(stream.conn.Close(), "Could not completely close the underlying TCP stream: %v")
}

/* ===============================================
			BEGIN PRIVATE PART
   =============================================== */

// nextTag returns the opening tag of the next XML element in the stream.
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

func writeToken(enc *xml.Encoder, token xml.Token) {
	utils.Successful(enc.EncodeToken(token), "Could not encode token %v")
	utils.Successful(enc.Flush(), "Could not flush after writing token %v")
}

// nextElement returns the opening tag of the next XML element and the complete element in []byte form.
func (stream *Stream) nextElement() (*xml.StartElement, []byte) {
	var temp struct {
		Inner []byte `xml:",innerxml"`
	}
	start, _ := stream.nextTag()
	end := start.End()
	utils.Successful(stream.decoder.DecodeElement(&temp, start), "Could not decode next xml element: %v")

	buffer := new(bytes.Buffer)
	enc := xml.NewEncoder(buffer)

	writeToken(enc, *start)
	buffer.Write(temp.Inner)
	writeToken(enc, end)

	return start, buffer.Bytes()
}
