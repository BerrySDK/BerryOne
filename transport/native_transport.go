package transport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/BerrySDK/berryone/events"
	"github.com/BerrySDK/berryone/native"
	"github.com/BerrySDK/berryone/protocol"
)

type NativeTransportOptions struct {
	AuthFolder string
	Browser    native.CompanionBrowser
	Config     protocol.WhatsAppWebConfig
}

type NativeTransport struct {
	mu        sync.Mutex
	sessionID string
	sink      EventSink
	store     *native.FileStore
	socket    *native.SocketClient
	options   NativeTransportOptions
}

func NewNativeTransport(options NativeTransportOptions) *NativeTransport {
	authFolder := options.AuthFolder
	if authFolder == "" {
		authFolder = ".auth"
	}
	if options.Browser == (native.CompanionBrowser{}) {
		options.Browser = native.DefaultCompanionBrowser
	}
	if options.Config == (protocol.WhatsAppWebConfig{}) {
		options.Config = protocol.DefaultWhatsAppWebConfig
	}

	return &NativeTransport{
		store:   native.NewFileStore(authFolder),
		socket:  native.NewSocketClient(options.Config),
		options: options,
	}
}

func (t *NativeTransport) SetEventSink(sink EventSink) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sink = sink
}

func (t *NativeTransport) Connect(ctx context.Context, sessionID string, auth *events.AuthOptions) error {
	t.mu.Lock()
	t.sessionID = sessionID
	t.mu.Unlock()

	record, _, err := t.store.Ensure(sessionID)
	if err != nil {
		return err
	}

	method := events.AuthMethodQR
	if auth != nil && auth.Method != "" {
		method = auth.Method
	}

	switch method {
	case events.AuthMethodQR, events.AuthMethodLink, events.AuthMethodPairingCode:
	default:
		return fmt.Errorf("unsupported auth method %q", method)
	}

	if method == events.AuthMethodPairingCode {
		code := ""
		if auth != nil {
			code = auth.CustomPairingCode
		}
		if code == "" {
			code = generatePairingCode()
		}
		t.emit(events.EventAuthPairingCode, struct {
			SessionID   string
			PhoneNumber string
			Code        string
		}{
			SessionID:   sessionID,
			PhoneNumber: auth.PhoneNumber,
			Code:        code,
		})
	}

	if err := t.socket.Connect(ctx); err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native websocket dial failed: %v", err))
		return err
	}

	ephemeralKey, err := native.GenerateX25519KeyPair()
	if err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native ephemeral key generation failed: %v", err))
		return err
	}

	noise := native.NewNoiseHandler(ephemeralKey.Private, ephemeralKey.Public, nil)
	hello := native.EncodeHandshakeClientHello(ephemeralKey.Public)
	frame := noise.EncodeFrame(hello)

	if err := t.socket.WriteFrame(2, frame); err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native client hello send failed: %v", err))
		return err
	}

	_, rawFrame, err := t.socket.ReadFrame()
	if err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native handshake read failed: %v", err))
		return err
	}

	decodedFrame, err := noise.DecodeHandshakeFrame(rawFrame)
	if err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native handshake frame decode failed: %v", err))
		return err
	}

	handshake, err := native.DecodeHandshakeMessage(decodedFrame)
	if err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native server hello decode failed: %v", err))
		return err
	}

	if handshake.ServerHello == nil {
		err = fmt.Errorf("berryone native runtime: server did not return serverHello in the initial handshake")
		t.emitProtocolError(sessionID, err.Error())
		return err
	}

	if _, err := noise.ProcessServerHello(*handshake.ServerHello, record.Credentials.NoiseKey); err != nil {
		t.emitProtocolError(sessionID, fmt.Sprintf("native server hello processing failed: %v", err))
		return err
	}

	t.emitProtocolError(sessionID, native.ErrRegistrationPayloadPending.Error())
	return native.ErrRegistrationPayloadPending
}

func (t *NativeTransport) Disconnect(_ context.Context, reason string) error {
	if err := t.socket.Close(); err != nil {
		return err
	}
	t.emit(events.EventConnectionClose, events.ConnectionState{
		SessionID:      t.sessionID,
		DisconnectedAt: time.Now(),
		Reason:         reason,
	})
	return nil
}

func (t *NativeTransport) Reconnect(ctx context.Context) error {
	return t.Connect(ctx, t.sessionID, &events.AuthOptions{Method: events.AuthMethodQR})
}

func (t *NativeTransport) Logout(ctx context.Context) error {
	if err := t.Disconnect(ctx, "logout"); err != nil {
		return err
	}
	return t.store.Remove(t.sessionID)
}

func (t *NativeTransport) SendMessage(_ context.Context, _ string, _ map[string]any, _ map[string]any) (events.OutgoingMessage, error) {
	return nil, errors.New("berryone native runtime: sendMessage is not implemented yet")
}

func (t *NativeTransport) emit(event events.EventName, payload any) {
	t.mu.Lock()
	sink := t.sink
	t.mu.Unlock()
	if sink != nil {
		sink(event, payload)
	}
}

func (t *NativeTransport) emitProtocolError(sessionID, message string) {
	t.emit(events.EventProtocolError, struct {
		SessionID string
		Error     string
	}{
		SessionID: sessionID,
		Error:     message,
	})
}

func generatePairingCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "BERRYONE"
	}
	return hex.EncodeToString(bytes)
}
