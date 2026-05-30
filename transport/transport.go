package transport

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/BerrySDK/berryone/events"
)

type EventSink func(event events.EventName, payload any)

type MessageContent struct {
	Text string
}

type Transport interface {
	SetEventSink(sink EventSink)
	Connect(ctx context.Context, sessionID string, auth *events.AuthOptions) error
	Disconnect(ctx context.Context, reason string) error
	Reconnect(ctx context.Context) error
	Logout(ctx context.Context) error
	SendMessage(ctx context.Context, to string, content map[string]any, options map[string]any) (events.OutgoingMessage, error)
}

type InMemoryTransport struct {
	mu        sync.Mutex
	sessionID string
	connected bool
	sink      EventSink
	counter   uint64
}

func NewInMemoryTransport() *InMemoryTransport {
	return &InMemoryTransport{}
}

func (t *InMemoryTransport) SetEventSink(sink EventSink) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sink = sink
}

func (t *InMemoryTransport) Connect(_ context.Context, sessionID string, auth *events.AuthOptions) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.sessionID = sessionID
	t.connected = true
	method := events.AuthMethodQR
	if auth != nil && auth.Method != "" {
		method = auth.Method
	}

	switch method {
	case events.AuthMethodQR:
		t.emitLocked(events.EventQR, "qr:berryone-demo")
		t.emitLocked(events.EventAuthQR, struct {
			SessionID string
			Value     string
		}{SessionID: sessionID, Value: "qr:berryone-demo"})
	case events.AuthMethodLink:
		t.emitLocked(events.EventAuthLink, struct {
			SessionID string
			Value     string
		}{SessionID: sessionID, Value: "link:berryone-demo"})
	case events.AuthMethodPairingCode:
		code := "BERRYONE-PAIR-1234"
		if auth != nil && auth.CustomPairingCode != "" {
			code = auth.CustomPairingCode
		}
		t.emitLocked(events.EventAuthPairingCode, struct {
			SessionID   string
			PhoneNumber string
			Code        string
		}{SessionID: sessionID, PhoneNumber: auth.PhoneNumber, Code: code})
	default:
		return fmt.Errorf("unsupported auth method %q", method)
	}

	t.emitLocked(events.EventConnectionOpen, events.ConnectionState{
		SessionID:   sessionID,
		ConnectedAt: time.Now(),
	})
	return nil
}

func (t *InMemoryTransport) Disconnect(_ context.Context, reason string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.connected {
		return nil
	}
	t.connected = false
	t.emitLocked(events.EventConnectionClose, events.ConnectionState{
		SessionID:      t.sessionID,
		DisconnectedAt: time.Now(),
		Reason:         reason,
	})
	return nil
}

func (t *InMemoryTransport) Reconnect(ctx context.Context) error {
	if err := t.Disconnect(ctx, "reconnect"); err != nil {
		return err
	}
	return t.Connect(ctx, t.sessionID, &events.AuthOptions{Method: events.AuthMethodQR})
}

func (t *InMemoryTransport) Logout(ctx context.Context) error {
	return t.Disconnect(ctx, "logout")
}

func (t *InMemoryTransport) SendMessage(_ context.Context, to string, content map[string]any, _ map[string]any) (events.OutgoingMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.connected {
		return nil, errors.New("transport is not connected")
	}

	t.counter++
	base := events.BaseMessage{
		ID:        fmt.Sprintf("berryone-%d", t.counter),
		To:        to,
		Timestamp: time.Now(),
		Ack:       events.AckPending,
		Type:      detectKind(content),
	}
	message := buildOutgoingMessage(base, content)

	t.emitLocked(events.EventMessageAck, events.MessageAck{
		MessageID: base.ID,
		RemoteJID: to,
		Ack:       events.AckSent,
		UpdatedAt: time.Now(),
	})
	t.emitLocked(events.EventMessageSent, message)
	return message, nil
}

func buildOutgoingMessage(base events.BaseMessage, content map[string]any) events.OutgoingMessage {
	switch base.Type {
	case "carousel":
		payload, _ := content["carousel"].(events.CarouselMessagePayload)
		return events.CarouselMessage{
			BaseMessage: base,
			Carousel:    payload,
		}
	case "image":
		payload, _ := content["image"].(*events.MediaPayload)
		if payload == nil {
			return events.ImageMessage{BaseMessage: base}
		}
		return events.ImageMessage{
			BaseMessage: base,
			Media:       *payload,
		}
	case "audio":
		payload, _ := content["audio"].(*events.MediaPayload)
		if payload == nil {
			return events.AudioMessage{BaseMessage: base}
		}
		return events.AudioMessage{
			BaseMessage: base,
			Media:       *payload,
		}
	case "document":
		payload, _ := content["document"].(*events.MediaPayload)
		if payload == nil {
			return events.DocumentMessage{BaseMessage: base}
		}
		return events.DocumentMessage{
			BaseMessage: base,
			Media:       *payload,
		}
	case "buttons":
		payload, _ := content["buttons"].(*events.ButtonsPayload)
		if payload == nil {
			return events.ButtonsMessage{BaseMessage: base}
		}
		return events.ButtonsMessage{
			BaseMessage: base,
			Buttons:     *payload,
		}
	case "list":
		payload, _ := content["list"].(*events.ListPayload)
		if payload == nil {
			return events.ListMessage{BaseMessage: base}
		}
		return events.ListMessage{
			BaseMessage: base,
			List:        *payload,
		}
	case "interactive":
		payload, _ := content["interactive"].(*events.InteractivePayload)
		if payload == nil {
			return events.InteractiveMessage{BaseMessage: base}
		}
		return events.InteractiveMessage{
			BaseMessage: base,
			Interactive: *payload,
		}
	case "reaction":
		payload, _ := content["reaction"].(*events.ReactionMessage)
		if payload == nil {
			return events.ReactionMessage{BaseMessage: base}
		}
		return events.ReactionMessage{
			BaseMessage:     base,
			Emoji:           payload.Emoji,
			TargetMessageID: payload.TargetMessageID,
		}
	case "location":
		payload, _ := content["location"].(*events.LocationPayload)
		if payload == nil {
			return events.LocationMessage{BaseMessage: base}
		}
		return events.LocationMessage{
			BaseMessage: base,
			Location:    *payload,
		}
	case "contact":
		payload, _ := content["contact"].(*events.ContactPayload)
		if payload == nil {
			return events.ContactMessage{BaseMessage: base}
		}
		return events.ContactMessage{
			BaseMessage: base,
			Contact:     *payload,
		}
	default:
		return events.TextMessage{
			BaseMessage: base,
			Text:        fmt.Sprintf("%v", content["text"]),
		}
	}
}

func detectKind(content map[string]any) string {
	switch {
	case content["carousel"] != nil:
		return "carousel"
	case content["image"] != nil:
		return "image"
	case content["audio"] != nil:
		return "audio"
	case content["document"] != nil:
		return "document"
	case content["buttons"] != nil:
		return "buttons"
	case content["list"] != nil:
		return "list"
	case content["interactive"] != nil:
		return "interactive"
	case content["reaction"] != nil:
		return "reaction"
	case content["location"] != nil:
		return "location"
	case content["contact"] != nil || content["contacts"] != nil:
		return "contact"
	default:
		return "text"
	}
}

func (t *InMemoryTransport) emitLocked(event events.EventName, payload any) {
	if t.sink == nil {
		return
	}
	t.sink(event, payload)
}
