package websocketapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type WSDialogMessage struct {
	ckafka.DialogMessage
	Type string `json:"type"`
}

type WSViewedMessage struct {
	ckafka.ViewedMessage
	Type string `json:"type"`
}

type ServerSettings struct {
	ClientWriterBufferSize int
}

type Server struct {
	logger                    *czap.Logger
	upgrader                  *websocket.Upgrader
	newMessagesKafkaReader    *kafka.Reader
	viewedMessagesKafkaReader *kafka.Reader
	clientsMutex              sync.Mutex
	clients                   map[uint64]*websocket.Conn
	settings                  ServerSettings
}

func NewServer(upgrader *websocket.Upgrader, newMessagesKafkaReader *kafka.Reader,
	viewedMessagesKafkaReader *kafka.Reader, logger *czap.Logger,
) *Server {
	server := &Server{
		upgrader:                  upgrader,
		clients:                   make(map[uint64]*websocket.Conn),
		newMessagesKafkaReader:    newMessagesKafkaReader,
		viewedMessagesKafkaReader: viewedMessagesKafkaReader,
		settings: ServerSettings{
			ClientWriterBufferSize: 5,
		},
		logger: logger,
	}
	go func() {
		server.listenNewMessages(context.TODO())
	}()
	go func() {
		server.listenViewedMessages(context.TODO())
	}()
	return server
}

func (ws *Server) listenNewMessages(ctx context.Context) {
	// listen kafka new message events & notify WS writer/reader
	for {
		message, err := ws.newMessagesKafkaReader.ReadMessage(ctx)
		ws.logger.DebugContext(ctx, "Got kafka message")
		if err != nil {
			ws.logger.ErrorContext(ctx, "Got kafka message error", zap.Error(err))
			continue
		}
		parsedMessage := ckafka.DialogMessage{}
		if err = json.Unmarshal(message.Value, &parsedMessage); err != nil {
			ws.logger.ErrorContext(ctx, "Can't parse kafka message", zap.Error(err))
			continue
		}

		func() {
			ws.clientsMutex.Lock()
			defer ws.clientsMutex.Unlock()

			wsMessage := WSDialogMessage{
				DialogMessage: parsedMessage,
				Type:          "new_message",
			}

			wsConn, ok := ws.clients[wsMessage.ReciverID.ID]
			if !ok {
				ws.logger.DebugContext(ctx, "Can not send message. Reciver not in online", zap.Uint64("reciverID", wsMessage.ReciverID.ID))
				return
			}

			rawMessage, err := json.Marshal(wsMessage)
			if err != nil {
				ws.logger.ErrorContext(ctx, "Error to parse struct", zap.Error(err))
				return
			}
			wsConn.WriteMessage(websocket.TextMessage, rawMessage)
		}()
	}
}

func (ws *Server) listenViewedMessages(ctx context.Context) {
	for {
		message, err := ws.viewedMessagesKafkaReader.ReadMessage(ctx)
		ws.logger.DebugContext(ctx, "Got kafka message")
		if err != nil {
			ws.logger.ErrorContext(ctx, "Got kafka message error", zap.Error(err))
			continue
		}
		parsedMessage := ckafka.ViewedMessage{}
		if err = json.Unmarshal(message.Value, &parsedMessage); err != nil {
			ws.logger.ErrorContext(ctx, "Can't parse kafka message", zap.Error(err))
			continue
		}

		func() {
			ws.clientsMutex.Lock()
			defer ws.clientsMutex.Unlock()

			wsMessage := WSViewedMessage{
				ViewedMessage: parsedMessage,
				Type:          "viewed",
			}

			wsConn, ok := ws.clients[wsMessage.ReciverID.ID]
			if !ok {
				ws.logger.DebugContext(
					ctx, "Can not send message. Reciver not in online",
					zap.Uint64("reciverID", wsMessage.ReciverID.ID),
				)
				return
			}

			rawMessage, err := json.Marshal(wsMessage)
			if err != nil {
				ws.logger.ErrorContext(ctx, "Error to parse struct", zap.Error(err))
				return
			}
			wsConn.WriteMessage(websocket.TextMessage, rawMessage)
		}()
	}
}

func (ws *Server) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := strconv.ParseUint(r.Header.Get("X-Account-ID"), 10, 64)
	ws.logger.DebugContext(ctx, "Connecting user with ID", zap.Uint64("userID", userID))
	if err != nil {
		ws.logger.ErrorContext(ctx, "Can't parse Authorization", zap.Error(err))
		return
	}
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.ErrorContext(ctx, "Can't upgrade WS connection", zap.Error(err))
		return
	}
	defer func() {
		ws.logger.DebugContext(ctx, "WS connection closed", zap.Uint64("userID", userID))
		conn.Close()
		ws.clientsMutex.Lock()
		defer ws.clientsMutex.Unlock()
		delete(ws.clients, userID)
	}()

	func() {
		ws.clientsMutex.Lock()
		defer ws.clientsMutex.Unlock()

		// Reacreate connection or fobidden?
		ws.clients[userID] = conn
	}()
	client := newClient(conn, ws.logger)
	client.listenAndWrite(context.TODO())
}
