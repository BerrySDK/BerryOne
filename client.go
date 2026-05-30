package berryone

import (
	"context"
	"errors"
	"fmt"

	"github.com/BerrySDK/berryone/auth"
	"github.com/BerrySDK/berryone/events"
	"github.com/BerrySDK/berryone/messages"
	"github.com/BerrySDK/berryone/socket"
	"github.com/BerrySDK/berryone/store"
	"github.com/BerrySDK/berryone/transport"
)

const maxCarouselCards = 10

type SendMessageContent struct {
	AI                 bool
	Text               string
	Cards              []CarouselCard
	CarouselCardType   CarouselCardType
	Image              *MediaPayload
	Audio              *MediaPayload
	Document           *MediaPayload
	Caption            string
	Mimetype           string
	FileName           string
	PTT                bool
	Footer             string
	ButtonsMessage     *ButtonsPayload
	List               *ListPayload
	InteractiveMessage *InteractivePayload
	React              *ReactionMessage
	Location           *LocationPayload
	Contacts           *struct {
		DisplayName string
		Contacts    []ContactPayload
	}
}

type SendRawOptions struct {
	Quoted              any
	Mentions            []string
	ContextInfo         map[string]any
	EphemeralExpiration int
	ForwardingScore     int
	StatusJIDList       []string
}

type Client struct {
	options  ClientOptions
	bus      *events.BerryEventBus
	sessions *auth.SessionManager
	store    *store.MemoryStore
	socket   *socket.BerrySocket
	lastQR   string
}

func NewClient(options ClientOptions) (*Client, error) {
	if options.SessionID == "" {
		return nil, ErrSessionIDRequired
	}
	if options.Transport == nil {
		options.Transport = transport.NewInMemoryTransport()
	}
	if options.SessionStore == nil {
		options.SessionStore = auth.NewMemorySessionStore()
	}

	bus := events.NewBerryEventBus()
	sock := socket.New(socket.Options{
		SessionID:            options.SessionID,
		ReconnectMaxAttempts: options.ReconnectMaxAttempts,
		ReconnectDelayMs:     int(options.ReconnectDelay.Milliseconds()),
		AuthFolder:           options.AuthFolder,
		Auth:                 options.Auth,
	}, bus, options.Transport)

	client := &Client{
		options:  options,
		bus:      bus,
		sessions: auth.NewSessionManager(options.SessionStore),
		store:    store.NewMemoryStore(),
		socket:   sock,
	}
	client.bindInternals()
	return client, nil
}

func NewBerryOne(options ClientOptions) (*Client, error) {
	return NewClient(options)
}

func (c *Client) On(event EventName, handler EventHandler) func() {
	return c.bus.On(event, handler)
}

func (c *Client) Once(event EventName, handler EventHandler) func() {
	return c.bus.Once(event, handler)
}

func (c *Client) Connect(ctx context.Context, authOptions *AuthOptions) error {
	if _, err := c.sessions.Get(c.options.SessionID); err != nil {
		return err
	}
	return c.socket.Connect(ctx, authOptions)
}

func (c *Client) ConnectWithLink(ctx context.Context) error {
	return c.Connect(ctx, &AuthOptions{Method: AuthMethodLink})
}

func (c *Client) ConnectWithQR(ctx context.Context) error {
	return c.Connect(ctx, &AuthOptions{Method: AuthMethodQR})
}

func (c *Client) ConnectWithPairingCode(ctx context.Context, phoneNumber, customPairingCode string) error {
	return c.Connect(ctx, &AuthOptions{
		Method:            AuthMethodPairingCode,
		PhoneNumber:       phoneNumber,
		CustomPairingCode: customPairingCode,
	})
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.socket.Disconnect(ctx, "manual")
}

func (c *Client) Reconnect(ctx context.Context) error {
	return c.socket.Reconnect(ctx)
}

func (c *Client) Logout(ctx context.Context) error {
	if err := c.socket.Logout(ctx); err != nil {
		return err
	}
	return c.sessions.Clear(c.options.SessionID)
}

func (c *Client) GetQRCode() string {
	return c.lastQR
}

func (c *Client) SendText(ctx context.Context, to, text string) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{Text: text}, SendRawOptions{})
}

func (c *Client) SendImage(ctx context.Context, to string, mediaPayload MediaPayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{Image: &mediaPayload}, SendRawOptions{})
}

func (c *Client) SendAudio(ctx context.Context, to string, mediaPayload MediaPayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{Audio: &mediaPayload}, SendRawOptions{})
}

func (c *Client) SendDocument(ctx context.Context, to string, mediaPayload MediaPayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{Document: &mediaPayload}, SendRawOptions{})
}

func (c *Client) SendButtons(ctx context.Context, to string, payload ButtonsPayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{ButtonsMessage: &payload}, SendRawOptions{})
}

func (c *Client) SendList(ctx context.Context, to string, payload ListPayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{List: &payload}, SendRawOptions{})
}

func (c *Client) SendCarousel(ctx context.Context, to string, payload CarouselMessagePayload) (OutgoingMessage, error) {
	return c.SendMessage(ctx, to, SendMessageContent{
		Text:             payload.Text,
		Footer:           payload.Footer,
		Cards:            payload.Cards,
		CarouselCardType: payload.CarouselCardType,
		AI:               payload.AI,
	}, SendRawOptions{})
}

func (c *Client) SendMessage(ctx context.Context, to string, content SendMessageContent, options SendRawOptions) (OutgoingMessage, error) {
	if err := validateJID(to); err != nil {
		return nil, err
	}
	if err := c.validateMessageContent(content); err != nil {
		return nil, err
	}
	rawContent, _, err := c.normalizeOutgoingContent(to, content)
	if err != nil {
		return nil, err
	}
	rawOptions := map[string]any{
		"quoted":              options.Quoted,
		"mentions":            options.Mentions,
		"contextInfo":         options.ContextInfo,
		"ephemeralExpiration": options.EphemeralExpiration,
		"forwardingScore":     options.ForwardingScore,
		"statusJidList":       options.StatusJIDList,
	}
	return c.socket.SendTransportMessage(ctx, to, rawContent, rawOptions)
}

func (c *Client) normalizeOutgoingContent(to string, content SendMessageContent) (map[string]any, OutgoingMessage, error) {
	switch {
	case len(content.Cards) > 0:
		carousel := CarouselMessagePayload{
			Text:             content.Text,
			Footer:           content.Footer,
			Cards:            content.Cards,
			CarouselCardType: content.CarouselCardType,
			AI:               content.AI,
		}
		msg := messages.CreateCarouselMessage(to, carousel)
		return map[string]any{"carousel": carousel}, msg, nil
	case content.Text != "":
		msg := messages.CreateTextMessage(to, content.Text)
		return map[string]any{"text": content.Text, "ai": content.AI}, msg, nil
	case content.Image != nil:
		msg := messages.CreateImageMessage(to, *content.Image)
		return map[string]any{"image": content.Image, "caption": content.Image.Caption, "ai": content.AI}, msg, nil
	case content.Audio != nil:
		msg := messages.CreateAudioMessage(to, *content.Audio)
		return map[string]any{"audio": content.Audio, "ai": content.AI}, msg, nil
	case content.Document != nil:
		msg := messages.CreateDocumentMessage(to, *content.Document)
		return map[string]any{"document": content.Document, "ai": content.AI}, msg, nil
	case content.ButtonsMessage != nil:
		msg := messages.CreateButtonsMessage(to, *content.ButtonsMessage)
		return map[string]any{"buttons": content.ButtonsMessage}, msg, nil
	case content.List != nil:
		msg := messages.CreateListMessage(to, *content.List)
		return map[string]any{"list": content.List}, msg, nil
	case content.React != nil:
		msg := messages.CreateReactionMessage(to, content.React.Emoji, content.React.TargetMessageID)
		return map[string]any{"reaction": content.React}, msg, nil
	case content.Location != nil:
		msg := messages.CreateLocationMessage(to, *content.Location)
		return map[string]any{"location": content.Location, "ai": content.AI}, msg, nil
	case content.Contacts != nil && len(content.Contacts.Contacts) > 0:
		msg := messages.CreateContactMessage(to, content.Contacts.Contacts[0])
		return map[string]any{"contacts": content.Contacts, "ai": content.AI}, msg, nil
	case content.InteractiveMessage != nil:
		base := BaseMessage{ID: "", To: to}
		msg := InteractiveMessage{BaseMessage: base, Interactive: *content.InteractiveMessage}
		return map[string]any{"interactive": content.InteractiveMessage}, msg, nil
	default:
		return nil, nil, errors.New("unsupported message content")
	}
}

func (c *Client) validateMessageContent(content SendMessageContent) error {
	if len(content.Cards) > maxCarouselCards {
		return fmt.Errorf("carousel payload supports at most %d cards", maxCarouselCards)
	}
	if len(content.Cards) > 0 {
		for index, card := range content.Cards {
			hasImage := card.Image != nil
			hasVideo := card.Video != nil
			if !hasImage && !hasVideo {
				return fmt.Errorf("carousel card %d must contain image or video", index+1)
			}
			if hasImage && hasVideo {
				return fmt.Errorf("carousel card %d cannot contain both image and video", index+1)
			}
		}
		return nil
	}
	switch {
	case content.Text != "":
		return nil
	case content.Image != nil:
		return nil
	case content.Audio != nil:
		return nil
	case content.Document != nil:
		return nil
	case content.ButtonsMessage != nil:
		return nil
	case content.List != nil:
		return nil
	case content.InteractiveMessage != nil:
		return nil
	case content.React != nil:
		return nil
	case content.Location != nil:
		return nil
	case content.Contacts != nil && len(content.Contacts.Contacts) > 0:
		return nil
	default:
		return errors.New("unsupported message content. use text, image, audio, document, buttonsMessage, list, cards, interactiveMessage, react, location or contacts")
	}
}

func validateJID(jid string) error {
	if jid == "" {
		return errors.New("invalid WhatsApp JID: empty")
	}
	for _, r := range jid {
		if r == '@' {
			return nil
		}
	}
	return fmt.Errorf("invalid WhatsApp JID %q: expected something like 5511999999999@s.whatsapp.net", jid)
}

func ValidateSendMessageContent(content SendMessageContent) error {
	client := &Client{}
	return client.validateMessageContent(content)
}

func (c *Client) bindInternals() {
	c.bus.On(events.EventQR, func(payload any) {
		qr, ok := payload.(string)
		if !ok {
			return
		}
		c.lastQR = qr
		_, _ = c.sessions.Update(c.options.SessionID, events.AuthStateSnapshot{QR: qr})
	})

	c.bus.On(events.EventAuthLink, func(payload any) {
		data, ok := payload.(struct {
			SessionID string
			Value     string
		})
		if !ok {
			return
		}
		_, _ = c.sessions.Update(c.options.SessionID, events.AuthStateSnapshot{
			AuthMethod: events.AuthMethodLink,
			LinkCode:   data.Value,
		})
	})

	c.bus.On(events.EventAuthQR, func(payload any) {
		data, ok := payload.(struct {
			SessionID string
			Value     string
		})
		if !ok {
			return
		}
		c.lastQR = data.Value
		_, _ = c.sessions.Update(c.options.SessionID, events.AuthStateSnapshot{
			AuthMethod: events.AuthMethodQR,
			QR:         data.Value,
		})
	})

	c.bus.On(events.EventAuthPairingCode, func(payload any) {
		data, ok := payload.(struct {
			SessionID   string
			PhoneNumber string
			Code        string
		})
		if !ok {
			return
		}
		_, _ = c.sessions.Update(c.options.SessionID, events.AuthStateSnapshot{
			AuthMethod:  events.AuthMethodPairingCode,
			PairingCode: data.Code,
		})
	})

	c.bus.On(events.EventConnectionOpen, func(_ any) {
		session, _ := c.sessions.Update(c.options.SessionID, events.AuthStateSnapshot{
			Registered: true,
		})
		c.bus.Emit(events.EventAuthSuccess, session)
	})

	c.bus.On(events.EventProtocolError, func(payload any) {
		data, ok := payload.(struct {
			SessionID string
			Error     string
		})
		if !ok {
			return
		}
		c.bus.Emit(events.EventAuthError, data)
	})

	c.bus.On(events.EventMessageReceived, func(payload any) {
		message, ok := payload.(events.IncomingMessage)
		if !ok {
			return
		}
		c.store.UpsertMessages(c.options.SessionID, []events.IncomingMessage{message})
	})

	c.bus.On(events.EventMessageAck, func(payload any) {
		ack, ok := payload.(events.MessageAck)
		if !ok {
			return
		}
		c.store.UpsertAck(c.options.SessionID, ack)
	})

	c.bus.On(events.EventChatsUpdate, func(payload any) {
		chats, ok := payload.([]events.ChatRecord)
		if !ok {
			return
		}
		c.store.UpsertChats(c.options.SessionID, chats)
	})

	c.bus.On(events.EventSyncContacts, func(payload any) {
		contacts, ok := payload.([]events.ContactRecord)
		if !ok {
			return
		}
		c.store.UpsertContacts(c.options.SessionID, contacts)
	})

	c.bus.On(events.EventSyncGroups, func(payload any) {
		groups, ok := payload.([]events.GroupRecord)
		if !ok {
			return
		}
		c.store.UpsertGroups(c.options.SessionID, groups)
	})
}
