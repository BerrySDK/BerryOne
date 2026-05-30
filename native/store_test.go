package native

import (
	"testing"
)

func TestFileStoreEnsureAndLoad(t *testing.T) {
	root := t.TempDir()
	store := NewFileStore(root)

	record, created, err := store.Ensure("session-a")
	if err != nil {
		t.Fatalf("ensure failed: %v", err)
	}
	if !created {
		t.Fatalf("expected first ensure to create record")
	}
	if record.SessionID != "session-a" {
		t.Fatalf("unexpected session id: %s", record.SessionID)
	}
	if len(record.Credentials.NoiseKey.Public) == 0 {
		t.Fatalf("expected generated noise key")
	}

	loaded, err := store.Load("session-a")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Credentials.ClientID != record.Credentials.ClientID {
		t.Fatalf("expected persisted client id")
	}
}
