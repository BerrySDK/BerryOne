package native

import "errors"

var (
	ErrHandshakeNotImplemented = errors.New("berryone native runtime: websocket handshake is not implemented yet")
	ErrSocketNotConnected      = errors.New("berryone native runtime: websocket is not connected")
)
