package protocol

import (
	"encoding/json"
	"errors"
)

type Frame struct {
	Stanza  string            `json:"stanza"`
	Attrs   map[string]string `json:"attrs"`
	Payload []byte            `json:"payload"`
}

type WhatsAppWebConfig struct {
	WebSocketURL string
	Origin       string
	UserAgent    string
}

type ProtocolError struct {
	Message string
}

func (e ProtocolError) Error() string {
	return e.Message
}

type BinaryFrameCodec struct{}

func (BinaryFrameCodec) Encode(frame Frame) ([]byte, error) {
	return json.Marshal(frame)
}

func (BinaryFrameCodec) Decode(buffer []byte) (Frame, error) {
	var frame Frame
	if err := json.Unmarshal(buffer, &frame); err != nil {
		return Frame{}, ProtocolError{Message: "unable to decode incoming frame as Berry envelope: " + err.Error()}
	}
	if frame.Attrs == nil {
		frame.Attrs = map[string]string{}
	}
	if frame.Stanza == "" {
		return Frame{}, errors.New("decoded frame is missing stanza")
	}
	return frame, nil
}

var DefaultWhatsAppWebConfig = WhatsAppWebConfig{
	WebSocketURL: "wss://web.whatsapp.com/ws/chat",
	Origin:       "https://web.whatsapp.com",
	UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
}
