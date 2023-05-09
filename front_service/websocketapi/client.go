package websocketapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ClientSettings struct {
	Addr string
}

type Client struct {
	jwtManager jwtmanager.Manager
	dialer     *websocket.Dialer
	settings   *ClientSettings
	logger     *czap.Logger
	upgrader   *websocket.Upgrader
}

func NewClient(jwtManager jwtmanager.Manager, dialer *websocket.Dialer,
	settings *ClientSettings, upgrader *websocket.Upgrader, logger *czap.Logger,
) *Client {
	return &Client{
		jwtManager: jwtManager,
		dialer:     dialer,
		settings:   settings,
		upgrader:   upgrader,
		logger:     logger,
	}
}

func (c *Client) ProxifyWithAuth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("token")
	headers := http.Header{}
	claims, err := c.jwtManager.VerifyToken(ctx, accessToken)
	if err != nil {
		c.logger.ErrorContext(ctx, "Can not auth in WS connection", zap.Error(err), zap.String("token", accessToken))
		return ErrAuth
	}

	headers.Set("X-Account-ID", strconv.FormatUint(claims.UserID, 10))
	headers.Set("X-User-Role", claims.UserRole)

	// Connect to commit
	cometConn, _, err := c.dialer.Dial(c.settings.Addr, headers)
	if err != nil {
		c.logger.ErrorContext(ctx, "Can not open WS connection with commit", zap.Error(err))
		return err
	}
	defer cometConn.Close()

	// Upgrade current conn
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.ErrorContext(ctx, "Can not upgrade WS connection", zap.Error(err))
		return err
	}
	defer conn.Close()

	// Copy messages
	for {
		messageType, message, err := cometConn.ReadMessage()
		c.logger.DebugContext(
			ctx,
			"WS message",
			zap.ByteString("message", message),
			zap.Int("message type", messageType),
		)
		if messageType == websocket.CloseMessage {
			return nil
		} else if err != nil {
			c.logger.ErrorContext(ctx, "Can not read WS message", zap.Error(err))
			return err
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			c.logger.ErrorContext(ctx, "Can not write WS message", zap.Error(err))
			return err
		}
	}
}
