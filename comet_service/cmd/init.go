package main

import (
	"net/http"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/segmentio/kafka-go"
)

func NewKafkaReaderFromConfig(config *cviper.CustomViper) *kafka.Reader {
	kafkaAddress := config.GetStringRequired("KAFKA_BROKER_ADDRESS")
	readerConfig := kafka.ReaderConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    config.GetStringRequired("KAFKA_NEW_MESSAGES_TOPIC"),
		GroupID:  config.GetStringRequired("KAFKA_NEW_MESSAGES_CONSUMER_GROUP_NAME"),
		MaxBytes: 10e4, // 10MB
	}
	return kafka.NewReader(readerConfig)
}

func NewHTTPServerFromConfig(config *cviper.CustomViper, handler http.Handler) http.Server {
	return http.Server{
		ReadTimeout:       config.GetSecondsDurationRequired("SERVICE_READ_TIMEOUT_SECONDS"),
		WriteTimeout:      config.GetSecondsDurationRequired("SERVICE_WRITE_TIMEOUT_SECONDS"),
		IdleTimeout:       config.GetSecondsDurationRequired("SERVICE_IDLE_TIMEOUT_SECONDS"),
		ReadHeaderTimeout: config.GetSecondsDurationRequired("SERVICE_READ_HEADER_TIMEOUT_SECONDS"),
		Addr:              config.GetStringRequired("SERVICE_ADDRESS"),
		Handler:           handler,
	}
}
