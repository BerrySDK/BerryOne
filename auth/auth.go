package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/BerrySDK/berryone/events"
)

type SessionStore interface {
	Load(sessionID string) (*events.AuthStateSnapshot, error)
	Save(sessionID string, snapshot events.AuthStateSnapshot) error
	Remove(sessionID string) error
}

type SessionManager struct {
	sessionStore SessionStore
}

func NewSessionManager(sessionStore SessionStore) *SessionManager {
	return &SessionManager{sessionStore: sessionStore}
}

func (m *SessionManager) Get(sessionID string) (events.AuthStateSnapshot, error) {
	snapshot, err := m.sessionStore.Load(sessionID)
	if err != nil {
		return events.AuthStateSnapshot{}, err
	}
	if snapshot != nil {
		return *snapshot, nil
	}

	created := events.AuthStateSnapshot{
		SessionID:  sessionID,
		Registered: false,
		ClientID:   randomID(),
	}
	if err := m.sessionStore.Save(sessionID, created); err != nil {
		return events.AuthStateSnapshot{}, err
	}
	return created, nil
}

func (m *SessionManager) Update(sessionID string, partial events.AuthStateSnapshot) (events.AuthStateSnapshot, error) {
	current, err := m.Get(sessionID)
	if err != nil {
		return events.AuthStateSnapshot{}, err
	}

	if partial.SessionID != "" {
		current.SessionID = partial.SessionID
	}
	current.Registered = partial.Registered || current.Registered
	if partial.ClientID != "" {
		current.ClientID = partial.ClientID
	}
	if partial.ServerToken != "" {
		current.ServerToken = partial.ServerToken
	}
	if partial.ClientToken != "" {
		current.ClientToken = partial.ClientToken
	}
	if partial.QR != "" {
		current.QR = partial.QR
	}
	if partial.LinkCode != "" {
		current.LinkCode = partial.LinkCode
	}
	if partial.PairingCode != "" {
		current.PairingCode = partial.PairingCode
	}
	if partial.AuthMethod != "" {
		current.AuthMethod = partial.AuthMethod
	}

	if err := m.sessionStore.Save(sessionID, current); err != nil {
		return events.AuthStateSnapshot{}, err
	}
	return current, nil
}

func (m *SessionManager) Clear(sessionID string) error {
	return m.sessionStore.Remove(sessionID)
}

func randomID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "berryone-client"
	}
	return hex.EncodeToString(bytes)
}
