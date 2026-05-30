package auth

import (
	"sync"

	"github.com/BerrySDK/berryone/events"
)

type MemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]events.AuthStateSnapshot
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		sessions: make(map[string]events.AuthStateSnapshot),
	}
}

func (s *MemorySessionStore) Load(sessionID string) (*events.AuthStateSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok := s.sessions[sessionID]
	if !ok {
		return nil, nil
	}
	copy := snapshot
	return &copy, nil
}

func (s *MemorySessionStore) Save(sessionID string, snapshot events.AuthStateSnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = snapshot
	return nil
}

func (s *MemorySessionStore) Remove(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
	return nil
}
