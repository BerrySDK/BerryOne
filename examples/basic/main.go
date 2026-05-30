package main

import (
	"context"
	"fmt"
	"log"

	berryone "github.com/BerrySDK/berryone"
)

func main() {
	client, err := berryone.NewBerryOne(berryone.ClientOptions{
		SessionID: "demo-session",
	})
	if err != nil {
		log.Fatal(err)
	}

	client.On(berryone.EventAuthQR, func(payload any) {
		fmt.Printf("received QR payload: %#v\n", payload)
	})

	client.On(berryone.EventAuthSuccess, func(payload any) {
		fmt.Printf("auth success: %#v\n", payload)
	})

	if err := client.ConnectWithQR(context.Background()); err != nil {
		log.Fatal(err)
	}

	response, err := client.SendText(
		context.Background(),
		"5511999999999@s.whatsapp.net",
		"Hello from BerryOne Go Edition",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("message response: %#v\n", response)
}
