package store

import (
	"sync"

	"github.com/BerrySDK/berryone/events"
)

type MemoryStore struct {
	mu       sync.RWMutex
	chats    map[string]map[string]events.ChatRecord
	contacts map[string]map[string]events.ContactRecord
	groups   map[string]map[string]events.GroupRecord
	messages map[string]map[string]events.IncomingMessage
	acks     map[string]map[string]events.MessageAck
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		chats:    make(map[string]map[string]events.ChatRecord),
		contacts: make(map[string]map[string]events.ContactRecord),
		groups:   make(map[string]map[string]events.GroupRecord),
		messages: make(map[string]map[string]events.IncomingMessage),
		acks:     make(map[string]map[string]events.MessageAck),
	}
}

func (s *MemoryStore) UpsertChats(sessionID string, chats []events.ChatRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.chats[sessionID] == nil {
		s.chats[sessionID] = map[string]events.ChatRecord{}
	}
	for _, chat := range chats {
		s.chats[sessionID][chat.ID] = chat
	}
}

func (s *MemoryStore) UpsertContacts(sessionID string, contacts []events.ContactRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.contacts[sessionID] == nil {
		s.contacts[sessionID] = map[string]events.ContactRecord{}
	}
	for _, contact := range contacts {
		s.contacts[sessionID][contact.ID] = contact
	}
}

func (s *MemoryStore) UpsertGroups(sessionID string, groups []events.GroupRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.groups[sessionID] == nil {
		s.groups[sessionID] = map[string]events.GroupRecord{}
	}
	for _, group := range groups {
		s.groups[sessionID][group.ID] = group
	}
}

func (s *MemoryStore) UpsertMessages(sessionID string, messages []events.IncomingMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.messages[sessionID] == nil {
		s.messages[sessionID] = map[string]events.IncomingMessage{}
	}
	for _, message := range messages {
		s.messages[sessionID][message.ID] = message
	}
}

func (s *MemoryStore) UpsertAck(sessionID string, ack events.MessageAck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.acks[sessionID] == nil {
		s.acks[sessionID] = map[string]events.MessageAck{}
	}
	s.acks[sessionID][ack.MessageID] = ack
}
