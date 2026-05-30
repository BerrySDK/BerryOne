package messages

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/BerrySDK/berryone/events"
)

func buildBase(kind string, to string) events.BaseMessage {
	return events.BaseMessage{
		ID:        randomID(),
		To:        to,
		Timestamp: time.Now(),
		Ack:       events.AckPending,
		Type:      kind,
	}
}

func CreateTextMessage(to, text string) events.TextMessage {
	return events.TextMessage{
		BaseMessage: buildBase("text", to),
		Text:        text,
	}
}

func CreateImageMessage(to string, media events.MediaPayload) events.ImageMessage {
	return events.ImageMessage{
		BaseMessage: buildBase("image", to),
		Media:       media,
	}
}

func CreateAudioMessage(to string, media events.MediaPayload) events.AudioMessage {
	return events.AudioMessage{
		BaseMessage: buildBase("audio", to),
		Media:       media,
	}
}

func CreateDocumentMessage(to string, media events.MediaPayload) events.DocumentMessage {
	return events.DocumentMessage{
		BaseMessage: buildBase("document", to),
		Media:       media,
	}
}

func CreateButtonsMessage(to string, buttons events.ButtonsPayload) events.ButtonsMessage {
	return events.ButtonsMessage{
		BaseMessage: buildBase("buttons", to),
		Buttons:     buttons,
	}
}

func CreateListMessage(to string, list events.ListPayload) events.ListMessage {
	return events.ListMessage{
		BaseMessage: buildBase("list", to),
		List:        list,
	}
}

func CreateCarouselMessage(to string, carousel events.CarouselMessagePayload) events.CarouselMessage {
	return events.CarouselMessage{
		BaseMessage: buildBase("carousel", to),
		Carousel:    carousel,
	}
}

func CreateReactionMessage(to, emoji, targetMessageID string) events.ReactionMessage {
	return events.ReactionMessage{
		BaseMessage:     buildBase("reaction", to),
		Emoji:           emoji,
		TargetMessageID: targetMessageID,
	}
}

func CreateLocationMessage(to string, location events.LocationPayload) events.LocationMessage {
	return events.LocationMessage{
		BaseMessage: buildBase("location", to),
		Location:    location,
	}
}

func CreateContactMessage(to string, contact events.ContactPayload) events.ContactMessage {
	return events.ContactMessage{
		BaseMessage: buildBase("contact", to),
		Contact:     contact,
	}
}

func randomID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(bytes)
}
