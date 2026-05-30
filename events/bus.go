package events

import "sync"

type EventName string

const (
	EventQR                     EventName = "qr"
	EventAuthLink               EventName = "auth.link"
	EventAuthQR                 EventName = "auth.qr"
	EventAuthPairingCode        EventName = "auth.pairing_code"
	EventConnectionOpen         EventName = "connection.open"
	EventConnectionClose        EventName = "connection.close"
	EventConnectionReconnecting EventName = "connection.reconnecting"
	EventAuthSuccess            EventName = "auth.success"
	EventAuthError              EventName = "auth.error"
	EventMessageReceived        EventName = "message.received"
	EventMessageSent            EventName = "message.sent"
	EventMessageAck             EventName = "message.ack"
	EventPresenceUpdate         EventName = "presence.update"
	EventChatsUpdate            EventName = "chats.update"
	EventSyncHistory            EventName = "sync.history"
	EventSyncContacts           EventName = "sync.contacts"
	EventSyncGroups             EventName = "sync.groups"
	EventSyncMessages           EventName = "sync.messages"
	EventRawFrame               EventName = "raw.frame"
	EventProtocolError          EventName = "protocol.error"
)

type EventHandler func(payload any)

type BerryEventBus struct {
	mu       sync.RWMutex
	handlers map[EventName]map[uint64]EventHandler
	nextID   uint64
}

func NewBerryEventBus() *BerryEventBus {
	return &BerryEventBus{
		handlers: make(map[EventName]map[uint64]EventHandler),
	}
}

func (b *BerryEventBus) On(event EventName, handler EventHandler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.nextID++
	id := b.nextID
	if b.handlers[event] == nil {
		b.handlers[event] = make(map[uint64]EventHandler)
	}
	b.handlers[event][id] = handler
	return func() { b.Off(event, id) }
}

func (b *BerryEventBus) Once(event EventName, handler EventHandler) func() {
	var off func()
	off = b.On(event, func(payload any) {
		off()
		handler(payload)
	})
	return off
}

func (b *BerryEventBus) Off(event EventName, id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if bucket := b.handlers[event]; bucket != nil {
		delete(bucket, id)
		if len(bucket) == 0 {
			delete(b.handlers, event)
		}
	}
}

func (b *BerryEventBus) Emit(event EventName, payload any) {
	b.mu.RLock()
	bucket := b.handlers[event]
	handlers := make([]EventHandler, 0, len(bucket))
	for _, handler := range bucket {
		handlers = append(handlers, handler)
	}
	b.mu.RUnlock()

	for _, handler := range handlers {
		handler(payload)
	}
}
