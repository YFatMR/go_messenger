package main

import (
	"strings"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/segmentio/kafka-go"
)

func NewKafkaReaderFromConfig(config *cviper.CustomViper) *kafka.Reader {
	kafkaAddress := config.GetStringRequired("KAFKA_READER_BROKER_ADDRESS")
	readerConfig := kafka.ReaderConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    config.GetStringRequired("KAFKA_READER_TOPIC"),
		GroupID:  config.GetStringRequired("KAFKA_READER_CONSUMER_GROUP_NAME"),
		MaxBytes: 10e4, // 10MB
	}
	return kafka.NewReader(readerConfig)
}

func NewKafkaWriterFromConfig(config *cviper.CustomViper) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(config.GetStringRequired("KAFKA_WRITER_BROKER_ADDRESS")),
		Balancer:     &kafka.LeastBytes{},
		Compression:  kafka.Snappy,
		Topic:        config.GetStringRequired("KAFKA_WRITER_TOPIC"),
		WriteTimeout: config.GetMillisecondsDurationRequired("KAFKA_WRITER_WRITE_TIMEOUT_MILLISECONDS"),
		ReadTimeout:  config.GetMillisecondsDurationRequired("KAFKA_WRITER_READ_TIMEOUT_MILLISECONDS"),
	}
}

func ExecutionCommandArgsFromConfig(config *cviper.CustomViper) []string {
	return strings.Split(config.GetStringRequired("EXECUTION_COMMAND_ARGS"), ";")
}
