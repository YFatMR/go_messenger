package websocketapi

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type client struct {
	logger *czap.Logger
	conn   *websocket.Conn
}

func newClient(conn *websocket.Conn, logger *czap.Logger) *client {
	return &client{
		conn:   conn,
		logger: logger,
	}
}

func (c *client) listenAndWrite(ctx context.Context) error {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if messageType == websocket.CloseMessage {
			return nil
		} else if err != nil {
			c.logger.ErrorContext(ctx, "Can not read WS message", zap.Error(err))
			return err
		}
		err = c.conn.WriteMessage(messageType, message)
		if err != nil {
			c.logger.ErrorContext(ctx, "Can not write WS message", zap.Error(err))
			return err
		}
	}
}
