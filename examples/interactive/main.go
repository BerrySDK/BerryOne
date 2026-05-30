package main

import (
	"context"
	"fmt"
	"log"

	berryone "github.com/BerrySDK/berryone"
)

func main() {
	client, err := berryone.NewBerryOne(berryone.ClientOptions{
		SessionID: "interactive-session",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := client.ConnectWithLink(context.Background()); err != nil {
		log.Fatal(err)
	}

	buttons := berryone.ButtonsPayload{
		Text:   "Choose an option",
		Footer: "Go Edition",
		Buttons: []berryone.ButtonRow{
			{ID: "docs", Title: "Open docs", Kind: berryone.ButtonKindReply},
			{ID: "support", Title: "Support", Kind: berryone.ButtonKindReply},
		},
	}

	buttonResponse, err := client.SendButtons(
		context.Background(),
		"5511999999999@s.whatsapp.net",
		buttons,
	)
	if err != nil {
		log.Fatal(err)
	}

	list := berryone.ListPayload{
		Title:      "BerryOne Menu",
		Text:       "Pick a workflow",
		ButtonText: "See options",
		Footer:     "BerryOne",
		Sections: []berryone.ListSection{
			{
				Title: "Main",
				Rows: []berryone.ListRow{
					{ID: "connect", Title: "Connect", Description: "Session onboarding"},
					{ID: "send", Title: "Send Message", Description: "Basic text flow"},
				},
			},
		},
	}

	listResponse, err := client.SendList(
		context.Background(),
		"5511999999999@s.whatsapp.net",
		list,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("buttons response: %#v\n", buttonResponse)
	fmt.Printf("list response: %#v\n", listResponse)
}
