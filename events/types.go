package events

import "time"

type AckStatus string

const (
	AckPending   AckStatus = "pending"
	AckSent      AckStatus = "sent"
	AckDelivered AckStatus = "delivered"
	AckRead      AckStatus = "read"
	AckFailed    AckStatus = "failed"
)

type ConnectionState struct {
	SessionID      string
	ConnectedAt    time.Time
	DisconnectedAt time.Time
	Reason         string
}

type AuthMethod string

const (
	AuthMethodLink        AuthMethod = "link"
	AuthMethodQR          AuthMethod = "qr"
	AuthMethodPairingCode AuthMethod = "pairing_code"
)

type AuthStateSnapshot struct {
	SessionID   string
	Registered  bool
	ClientID    string
	ServerToken string
	ClientToken string
	QR          string
	LinkCode    string
	PairingCode string
	AuthMethod  AuthMethod
}

type AuthOptions struct {
	Method            AuthMethod
	PhoneNumber       string
	CustomPairingCode string
}

type ContactRecord struct {
	ID        string
	Name      string
	PushName  string
	ShortName string
}

type ChatRecord struct {
	ID            string
	Name          string
	UnreadCount   int
	LastMessageAt time.Time
}

type GroupRecord struct {
	ID           string
	Subject      string
	Participants []string
}

type PresenceStatus string

const (
	PresenceAvailable   PresenceStatus = "available"
	PresenceComposing   PresenceStatus = "composing"
	PresenceRecording   PresenceStatus = "recording"
	PresencePaused      PresenceStatus = "paused"
	PresenceUnavailable PresenceStatus = "unavailable"
)

type PresenceRecord struct {
	ID         string
	Status     PresenceStatus
	LastSeenAt time.Time
}

type LocationPayload struct {
	Latitude  float64
	Longitude float64
	Name      string
	Address   string
}

type ContactPayload struct {
	DisplayName string
	VCard       string
}

type MediaPayload struct {
	URL      string
	Path     string
	Buffer   []byte
	FileName string
	Mimetype string
	Caption  string
}

type ButtonKind string

const (
	ButtonKindReply      ButtonKind = "reply"
	ButtonKindQuickReply ButtonKind = "quick_reply"
	ButtonKindCopyCode   ButtonKind = "copy_code"
	ButtonKindCTAURL     ButtonKind = "cta_url"
)

type ButtonRow struct {
	ID               string
	Title            string
	Kind             ButtonKind
	Code             string
	URL              string
	NativeFlowName   string
	ButtonParamsJSON string
}

type InteractiveHeader struct {
	Title              string
	Subtitle           string
	HasMediaAttachment bool
}

type InteractiveBody struct {
	Text string
}

type InteractiveFooter struct {
	Text string
}

type InteractiveNativeButton struct {
	Name             string
	ButtonParamsJSON string
}

type InteractiveNativeFlowPayload struct {
	Buttons           []InteractiveNativeButton
	MessageParamsJSON string
	MessageVersion    int
}

type InteractivePayload struct {
	Header            *InteractiveHeader
	Body              InteractiveBody
	Footer            *InteractiveFooter
	NativeFlowMessage *InteractiveNativeFlowPayload
}

type ButtonsPayload struct {
	Text    string
	Footer  string
	Buttons []ButtonRow
}

type CarouselCardType string

const (
	CarouselCardTypeImage CarouselCardType = "image"
	CarouselCardTypeVideo CarouselCardType = "video"
	CarouselCardTypeMixed CarouselCardType = "mixed"
)

type CarouselButton struct {
	ID               string
	Title            string
	Kind             ButtonKind
	Code             string
	URL              string
	NativeFlowName   string
	ButtonParamsJSON string
	Name             string
}

type CarouselCard struct {
	Title   string
	Body    string
	Footer  string
	Image   *MediaPayload
	Video   *MediaPayload
	Buttons []CarouselButton
}

type CarouselMessagePayload struct {
	Text             string
	Footer           string
	Cards            []CarouselCard
	CarouselCardType CarouselCardType
	AI               bool
}

type ListRow struct {
	ID          string
	Title       string
	Description string
}

type ListSection struct {
	Title string
	Rows  []ListRow
}

type ListPayload struct {
	Title      string
	Text       string
	Footer     string
	ButtonText string
	Sections   []ListSection
}

type BaseMessage struct {
	ID                  string
	To                  string
	ChatID              string
	RemoteJID           string
	From                string
	Timestamp           time.Time
	Ack                 AckStatus
	ButtonID            string
	SelectedButtonID    string
	RawButtonParamsJSON string
	Type                string
}

type TextMessage struct {
	BaseMessage
	Text string
}

type ImageMessage struct {
	BaseMessage
	Media MediaPayload
}

type AudioMessage struct {
	BaseMessage
	Media MediaPayload
}

type DocumentMessage struct {
	BaseMessage
	Media MediaPayload
}

type ButtonsMessage struct {
	BaseMessage
	Buttons ButtonsPayload
}

type ListMessage struct {
	BaseMessage
	List ListPayload
}

type CarouselMessage struct {
	BaseMessage
	Carousel CarouselMessagePayload
}

type InteractiveMessage struct {
	BaseMessage
	Interactive InteractivePayload
}

type ReactionMessage struct {
	BaseMessage
	Emoji           string
	TargetMessageID string
}

type LocationMessage struct {
	BaseMessage
	Location LocationPayload
}

type ContactMessage struct {
	BaseMessage
	Contact ContactPayload
}

type OutgoingMessage interface {
	GetBase() BaseMessage
}

func (m TextMessage) GetBase() BaseMessage        { return m.BaseMessage }
func (m ImageMessage) GetBase() BaseMessage       { return m.BaseMessage }
func (m AudioMessage) GetBase() BaseMessage       { return m.BaseMessage }
func (m DocumentMessage) GetBase() BaseMessage    { return m.BaseMessage }
func (m ButtonsMessage) GetBase() BaseMessage     { return m.BaseMessage }
func (m ListMessage) GetBase() BaseMessage        { return m.BaseMessage }
func (m CarouselMessage) GetBase() BaseMessage    { return m.BaseMessage }
func (m InteractiveMessage) GetBase() BaseMessage { return m.BaseMessage }
func (m ReactionMessage) GetBase() BaseMessage    { return m.BaseMessage }
func (m LocationMessage) GetBase() BaseMessage    { return m.BaseMessage }
func (m ContactMessage) GetBase() BaseMessage     { return m.BaseMessage }

type IncomingMessage struct {
	BaseMessage
	Text string
}

type MessageAck struct {
	MessageID string
	RemoteJID string
	Ack       AckStatus
	UpdatedAt time.Time
}

type SyncBundle struct {
	Contacts []ContactRecord
	Chats    []ChatRecord
	Groups   []GroupRecord
	Messages []IncomingMessage
}
