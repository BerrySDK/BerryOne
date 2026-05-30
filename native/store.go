package native

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/BerrySDK/berryone/events"
)

type StateRecord struct {
	SessionID   string                   `json:"session_id"`
	Credentials AuthCredentials          `json:"credentials"`
	Snapshot    events.AuthStateSnapshot `json:"snapshot"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type FileStore struct {
	RootDir string
}

func NewFileStore(rootDir string) *FileStore {
	return &FileStore{RootDir: rootDir}
}

func (s *FileStore) Ensure(sessionID string) (StateRecord, bool, error) {
	record, err := s.Load(sessionID)
	if err == nil && record != nil {
		return *record, false, nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return StateRecord{}, false, err
	}

	creds, err := GenerateAuthCredentials()
	if err != nil {
		return StateRecord{}, false, err
	}

	now := time.Now().UTC()
	record = &StateRecord{
		SessionID:   sessionID,
		Credentials: creds,
		Snapshot: events.AuthStateSnapshot{
			SessionID:  sessionID,
			Registered: false,
			ClientID:   creds.ClientID,
		},
		UpdatedAt: now,
	}

	if err := s.Save(*record); err != nil {
		return StateRecord{}, false, err
	}

	return *record, true, nil
}

func (s *FileStore) Load(sessionID string) (*StateRecord, error) {
	path := s.statePath(sessionID)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var record StateRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	if record.SessionID == "" {
		record.SessionID = sessionID
	}

	return &record, nil
}

func (s *FileStore) Save(record StateRecord) error {
	record.UpdatedAt = time.Now().UTC()
	if err := os.MkdirAll(filepath.Dir(s.statePath(record.SessionID)), 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.statePath(record.SessionID), payload, 0o600)
}

func (s *FileStore) Remove(sessionID string) error {
	path := filepath.Dir(s.statePath(sessionID))
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) statePath(sessionID string) string {
	return filepath.Join(s.RootDir, sessionID, "state.json")
}
