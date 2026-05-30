package berryone

import (
	"time"

	"github.com/BerrySDK/berryone/auth"
	"github.com/BerrySDK/berryone/events"
	"github.com/BerrySDK/berryone/media"
	"github.com/BerrySDK/berryone/protocol"
	"github.com/BerrySDK/berryone/store"
	"github.com/BerrySDK/berryone/transport"
)

type AuthMethod = events.AuthMethod
type AuthOptions = events.AuthOptions
type AckStatus = events.AckStatus

const (
	AuthMethodLink        = events.AuthMethodLink
	AuthMethodQR          = events.AuthMethodQR
	AuthMethodPairingCode = events.AuthMethodPairingCode

	AckPending   = events.AckPending
	AckSent      = events.AckSent
	AckDelivered = events.AckDelivered
	AckRead      = events.AckRead
	AckFailed    = events.AckFailed
)

type (
	ConnectionState              = events.ConnectionState
	AuthStateSnapshot            = events.AuthStateSnapshot
	ContactRecord                = events.ContactRecord
	ChatRecord                   = events.ChatRecord
	GroupRecord                  = events.GroupRecord
	PresenceRecord               = events.PresenceRecord
	LocationPayload              = events.LocationPayload
	ContactPayload               = events.ContactPayload
	MediaPayload                 = events.MediaPayload
	ButtonKind                   = events.ButtonKind
	ButtonRow                    = events.ButtonRow
	InteractiveHeader            = events.InteractiveHeader
	InteractiveBody              = events.InteractiveBody
	InteractiveFooter            = events.InteractiveFooter
	InteractiveNativeButton      = events.InteractiveNativeButton
	InteractiveNativeFlowPayload = events.InteractiveNativeFlowPayload
	InteractivePayload           = events.InteractivePayload
	ButtonsPayload               = events.ButtonsPayload
	CarouselCardType             = events.CarouselCardType
	CarouselButton               = events.CarouselButton
	CarouselCard                 = events.CarouselCard
	CarouselMessagePayload       = events.CarouselMessagePayload
	ListRow                      = events.ListRow
	ListSection                  = events.ListSection
	ListPayload                  = events.ListPayload
	BaseMessage                  = events.BaseMessage
	TextMessage                  = events.TextMessage
	ImageMessage                 = events.ImageMessage
	AudioMessage                 = events.AudioMessage
	DocumentMessage              = events.DocumentMessage
	ButtonsMessage               = events.ButtonsMessage
	ListMessage                  = events.ListMessage
	CarouselMessage              = events.CarouselMessage
	InteractiveMessage           = events.InteractiveMessage
	ReactionMessage              = events.ReactionMessage
	LocationMessage              = events.LocationMessage
	ContactMessage               = events.ContactMessage
	OutgoingMessage              = events.OutgoingMessage
	IncomingMessage              = events.IncomingMessage
	MessageAck                   = events.MessageAck
	SyncBundle                   = events.SyncBundle
	EventName                    = events.EventName
	EventHandler                 = events.EventHandler
	BerryEventBus                = events.BerryEventBus
	SessionManager               = auth.SessionManager
	SessionStore                 = auth.SessionStore
	MemorySessionStore           = auth.MemorySessionStore
	FileSessionStore             = auth.FileSessionStore
	LoadedMedia                  = media.LoadedMedia
	MediaManager                 = media.Manager
	BinaryFrameCodec             = protocol.BinaryFrameCodec
	ProtocolFrame                = protocol.Frame
	WhatsAppWebConfig            = protocol.WhatsAppWebConfig
	Transport                    = transport.Transport
	InMemoryTransport            = transport.InMemoryTransport
)

const (
	EventQR                     = events.EventQR
	EventAuthLink               = events.EventAuthLink
	EventAuthQR                 = events.EventAuthQR
	EventAuthPairingCode        = events.EventAuthPairingCode
	EventConnectionOpen         = events.EventConnectionOpen
	EventConnectionClose        = events.EventConnectionClose
	EventConnectionReconnecting = events.EventConnectionReconnecting
	EventAuthSuccess            = events.EventAuthSuccess
	EventAuthError              = events.EventAuthError
	EventMessageReceived        = events.EventMessageReceived
	EventMessageSent            = events.EventMessageSent
	EventMessageAck             = events.EventMessageAck
	EventPresenceUpdate         = events.EventPresenceUpdate
	EventChatsUpdate            = events.EventChatsUpdate
	EventSyncHistory            = events.EventSyncHistory
	EventSyncContacts           = events.EventSyncContacts
	EventSyncGroups             = events.EventSyncGroups
	EventSyncMessages           = events.EventSyncMessages
	EventRawFrame               = events.EventRawFrame
	EventProtocolError          = events.EventProtocolError
)

const (
	ButtonKindReply      = events.ButtonKindReply
	ButtonKindQuickReply = events.ButtonKindQuickReply
	ButtonKindCopyCode   = events.ButtonKindCopyCode
	ButtonKindCTAURL     = events.ButtonKindCTAURL

	CarouselCardTypeImage = events.CarouselCardTypeImage
	CarouselCardTypeVideo = events.CarouselCardTypeVideo
	CarouselCardTypeMixed = events.CarouselCardTypeMixed
)

type ClientOptions struct {
	SessionID            string
	DatabasePath         string
	AuthFolder           string
	Auth                 *AuthOptions
	ReconnectMaxAttempts int
	ReconnectDelay       time.Duration
	PrintQRInTerminal    bool
	QRSmall              bool
	Transport            Transport
	SessionStore         SessionStore
}

func NewEventBus() *BerryEventBus {
	return events.NewBerryEventBus()
}

func NewMemoryStore() *store.MemoryStore {
	return store.NewMemoryStore()
}

func NewMediaManager() *MediaManager {
	return media.NewManager()
}

func NewMemorySessionStore() *MemorySessionStore {
	return auth.NewMemorySessionStore()
}

func NewFileSessionStore(path string) *FileSessionStore {
	return auth.NewFileSessionStore(path)
}

func NewInMemoryTransport() *InMemoryTransport {
	return transport.NewInMemoryTransport()
}

var DefaultWhatsAppWebConfig = protocol.DefaultWhatsAppWebConfig
