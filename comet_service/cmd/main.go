package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/YFatMR/go_messenger/comet_service/internal/websocketapi"
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func main() {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	// Init environment vars

	// Init logger
	logger, err := czap.FromConfig(config)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Debug("Starting server...")

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic!", zap.Any("msg", r))
		}
		panic("Panic")
	}()

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	newMessageKafkaReader := NewKafkaReaderFromConfig(config)
	websocketServer := websocketapi.NewServer(upgrader, newMessageKafkaReader, logger)

	router := mux.NewRouter()
	router.HandleFunc("/", websocketServer.Handle)

	service := NewHTTPServerFromConfig(config, router)
	logger.Info(
		"Starting to register WS comet server",
		zap.String("WS comet server address", service.Addr),
	)
	if err := service.ListenAndServe(); err != nil {
		panic(err)
	}
}
