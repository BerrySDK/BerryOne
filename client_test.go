package berryone

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/BerrySDK/berryone/auth"
)

func TestConnectWithQRPersistsSessionData(t *testing.T) {
	store := NewMemorySessionStore()
	client, err := NewClient(ClientOptions{
		SessionID:    "session-a",
		SessionStore: store,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ConnectWithQR(context.Background()); err != nil {
		t.Fatalf("ConnectWithQR() error = %v", err)
	}

	manager := auth.NewSessionManager(store)
	session, err := manager.Get("session-a")
	if err != nil {
		t.Fatalf("manager.Get() error = %v", err)
	}

	if !session.Registered {
		t.Fatalf("expected session to be registered")
	}
	if session.AuthMethod != AuthMethodQR {
		t.Fatalf("expected QR auth method, got %q", session.AuthMethod)
	}
	if session.QR == "" {
		t.Fatalf("expected QR to be stored")
	}
}

func TestSendTextReturnsOutgoingMessage(t *testing.T) {
	client, err := NewClient(ClientOptions{
		SessionID: "session-b",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if err := client.ConnectWithQR(context.Background()); err != nil {
		t.Fatalf("ConnectWithQR() error = %v", err)
	}

	message, err := client.SendText(context.Background(), "5511999999999@s.whatsapp.net", "hello")
	if err != nil {
		t.Fatalf("SendText() error = %v", err)
	}

	base := message.GetBase()
	if base.ID == "" {
		t.Fatalf("expected message ID")
	}
	if base.Type != "text" {
		t.Fatalf("expected text kind, got %q", base.Type)
	}
}

func TestCarouselValidation(t *testing.T) {
	err := ValidateSendMessageContent(SendMessageContent{
		Cards: []CarouselCard{
			{Title: "bad card"},
		},
		CarouselCardType: CarouselCardTypeImage,
	})
	if err == nil {
		t.Fatalf("expected carousel validation error")
	}
}

func TestFileSessionStoreRoundTrip(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "sessions.json")
	store := NewFileSessionStore(storePath)
	manager := auth.NewSessionManager(store)

	record, err := manager.Update("session-c", AuthStateSnapshot{
		AuthMethod: AuthMethodLink,
		QR:         "qr:test",
	})
	if err != nil {
		t.Fatalf("manager.Update() error = %v", err)
	}

	if record.AuthMethod != AuthMethodLink {
		t.Fatalf("expected auth method to be persisted")
	}

	loaded, err := manager.Get("session-c")
	if err != nil {
		t.Fatalf("manager.Get() error = %v", err)
	}

	if loaded.QR != "qr:test" {
		t.Fatalf("expected qr %q, got %q", "qr:test", loaded.QR)
	}
}
