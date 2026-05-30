# BerryOne

BerryOne is the Go edition of the BerryProtocol developer experience.

It is designed to feel familiar to developers who already use the TypeScript SDK:

- client lifecycle helpers
- QR, link, and pairing-code connection flows
- event subscriptions
- session persistence
- text, media, buttons, lists, carousels, reactions, locations, and contacts
- transport abstraction ready for a real WhatsApp Web runtime in Go

## Status

BerryOne already compiles, has tests, ships examples, and exposes a stable public API.

What is fully working today:

- Go module structure
- public SDK surface
- message payload validation
- event bus
- in-memory transport for local development
- in-memory and file-based session stores
- examples and tests

What is still pending for full parity with the TypeScript BerryProtocol runtime:

- real WhatsApp Web transport
- QR rendering against a live session
- auth/session cryptography
- media upload pipeline
- sync, retries, app-state, and reconnection logic against the real network

## Install

```bash
go get github.com/BerrySDK/berryone
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	berryone "github.com/BerrySDK/berryone"
)

func main() {
	client, err := berryone.NewBerryOne(berryone.ClientOptions{
		SessionID: "default",
	})
	if err != nil {
		log.Fatal(err)
	}

	client.On(berryone.EventAuthQR, func(payload any) {
		fmt.Printf("qr payload: %#v\n", payload)
	})

	client.On(berryone.EventAuthSuccess, func(payload any) {
		fmt.Printf("auth success: %#v\n", payload)
	})

	if err := client.ConnectWithQR(context.Background()); err != nil {
		log.Fatal(err)
	}

	message, err := client.SendText(
		context.Background(),
		"5511999999999@s.whatsapp.net",
		"Hello from BerryOne",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("message sent: %#v\n", message)
}
```

## Examples

- [C:\Users\felip\BerryOne\examples\basic\main.go](C:\Users\felip\BerryOne\examples\basic\main.go)
- [C:\Users\felip\BerryOne\examples\interactive\main.go](C:\Users\felip\BerryOne\examples\interactive\main.go)

## Project Layout

- `auth/`: session stores and session manager
- `events/`: event bus, payloads, message and sync types
- `media/`: media loading helpers
- `messages/`: outgoing message builders
- `protocol/`: frame codec and protocol config
- `socket/`: socket facade over the transport layer
- `store/`: in-memory state store for chats, contacts, groups, messages and acks
- `transport/`: transport contracts and in-memory runtime
- `client.go`: public high-level SDK facade
- `exports.go`: package-level aliases and convenience constructors
- `examples/`: runnable examples
- `client_test.go`: baseline integration tests for the public API

## Development

Run tests:

```bash
go test ./...
```

Run the basic example:

```bash
go run ./examples/basic
```

## Roadmap

- real socket/auth runtime in Go
- persistent message and chat store
- media loader and uploader
- group and presence operations
- webhook/API layer for BerryOne services

## License

Apache-2.0
