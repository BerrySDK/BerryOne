package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	berryone "github.com/BerrySDK/berryone"
)

func main() {
	client, err := berryone.NewBerryOne(berryone.ClientOptions{
		SessionID:         "native-session",
		AuthFolder:        ".auth",
		PrintQRInTerminal: true,
		QRSmall:           true,
		Transport: berryone.NewNativeTransport(berryone.NativeTransportOptions{
			AuthFolder: ".auth",
		}),
		ReconnectDelay: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	client.On(berryone.EventQR, func(payload any) {
		fmt.Println("QR generated. Scan it with WhatsApp Linked Devices.")
	})

	client.On(berryone.EventProtocolError, func(payload any) {
		fmt.Printf("protocol error: %#v\n", payload)
	})

	err = client.ConnectWithQR(context.Background())
	if errors.Is(err, berryone.ErrNativeHandshakeNotImplemented) {
		fmt.Println("native status: QR rendering works, but the real WhatsApp handshake is still being ported.")
		fmt.Println("The process will stay open so you can inspect the QR. Press Enter to close.")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}
	if errors.Is(err, berryone.ErrNativeRegistrationPending) {
		fmt.Println("native status: the real websocket handshake reached the server hello stage.")
		fmt.Println("next step pending: port the registration payload so WhatsApp can return a real pairing QR.")
		return
	}
	fmt.Println("connect result:", err)
}
