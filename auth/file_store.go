package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/BerrySDK/berryone/events"
)

type FileSessionStore struct {
	mu   sync.Mutex
	path string
}

func NewFileSessionStore(path string) *FileSessionStore {
	return &FileSessionStore{path: path}
}

func (s *FileSessionStore) Load(sessionID string) (*events.AuthStateSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	records, err := s.readAll()
	if err != nil {
		return nil, err
	}
	snapshot, ok := records[sessionID]
	if !ok {
		return nil, nil
	}
	copy := snapshot
	return &copy, nil
}

func (s *FileSessionStore) Save(sessionID string, snapshot events.AuthStateSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	records, err := s.readAll()
	if err != nil {
		return err
	}
	records[sessionID] = snapshot
	return s.writeAll(records)
}

func (s *FileSessionStore) Remove(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	records, err := s.readAll()
	if err != nil {
		return err
	}
	delete(records, sessionID)
	return s.writeAll(records)
}

func (s *FileSessionStore) readAll() (map[string]events.AuthStateSnapshot, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]events.AuthStateSnapshot{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return map[string]events.AuthStateSnapshot{}, nil
	}
	records := map[string]events.AuthStateSnapshot{}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (s *FileSessionStore) writeAll(records map[string]events.AuthStateSnapshot) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
