package native

import (
	"context"
	"net/http"
	"time"

	"github.com/BerrySDK/berryone/protocol"
	"github.com/gorilla/websocket"
)

type SocketClient struct {
	Config protocol.WhatsAppWebConfig
	conn   *websocket.Conn
}

func NewSocketClient(config protocol.WhatsAppWebConfig) *SocketClient {
	return &SocketClient{Config: config}
}

func (c *SocketClient) Connect(ctx context.Context) error {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 20 * time.Second,
	}

	header := http.Header{}
	header.Set("Origin", c.Config.Origin)
	header.Set("User-Agent", c.Config.UserAgent)

	conn, _, err := dialer.DialContext(ctx, c.Config.WebSocketURL, header)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *SocketClient) Close() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *SocketClient) ReadFrame() (int, []byte, error) {
	if c.conn == nil {
		return 0, nil, ErrSocketNotConnected
	}
	return c.conn.ReadMessage()
}

func (c *SocketClient) WriteFrame(messageType int, payload []byte) error {
	if c.conn == nil {
		return ErrSocketNotConnected
	}
	return c.conn.WriteMessage(messageType, payload)
}
