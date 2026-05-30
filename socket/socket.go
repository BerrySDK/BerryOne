package socket

import (
	"context"
	"errors"
	"strings"

	"github.com/BerrySDK/berryone/events"
	"github.com/BerrySDK/berryone/transport"
)

type Options struct {
	SessionID            string
	ReconnectMaxAttempts int
	ReconnectDelayMs     int
	AuthFolder           string
	Auth                 *events.AuthOptions
}

type BerrySocket struct {
	options Options
	bus     *events.BerryEventBus
	engine  transport.Transport
	auth    events.AuthOptions
}

func New(options Options, bus *events.BerryEventBus, engine transport.Transport) *BerrySocket {
	if engine == nil {
		engine = transport.NewInMemoryTransport()
	}
	socket := &BerrySocket{
		options: options,
		bus:     bus,
		engine:  engine,
		auth:    events.AuthOptions{Method: events.AuthMethodLink},
	}
	if options.Auth != nil {
		socket.auth = *options.Auth
	}
	engine.SetEventSink(func(event events.EventName, payload any) {
		bus.Emit(event, payload)
	})
	return socket
}

func (s *BerrySocket) SetAuth(auth events.AuthOptions) {
	if auth.Method == "" {
		auth.Method = events.AuthMethodLink
	}
	if auth.Method == events.AuthMethodPairingCode {
		auth.PhoneNumber = normalizePhoneNumber(auth.PhoneNumber)
	}
	s.auth = auth
}

func (s *BerrySocket) Connect(ctx context.Context, auth *events.AuthOptions) error {
	if auth != nil {
		s.SetAuth(*auth)
	}
	return s.engine.Connect(ctx, s.options.SessionID, &s.auth)
}

func (s *BerrySocket) Disconnect(ctx context.Context, reason string) error {
	return s.engine.Disconnect(ctx, reason)
}

func (s *BerrySocket) Reconnect(ctx context.Context) error {
	return s.engine.Reconnect(ctx)
}

func (s *BerrySocket) Logout(ctx context.Context) error {
	return s.engine.Logout(ctx)
}

func (s *BerrySocket) SendTransportMessage(ctx context.Context, to string, content map[string]any, options map[string]any) (events.OutgoingMessage, error) {
	if !strings.Contains(to, "@") {
		return nil, errors.New("invalid WhatsApp JID")
	}
	return s.engine.SendMessage(ctx, to, content, options)
}

func normalizePhoneNumber(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
