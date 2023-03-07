package test

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/segmentio/kafka-go"
)

type KafkaLimitations struct {
	ReadTimeout time.Duration
}

type KafkaClient struct {
	Reader      *kafka.Reader
	Limitations KafkaLimitations
}

func (k *KafkaClient) WaitMessage(ctx context.Context) (
	kafka.Message, error,
) {
	ctx, cancel := context.WithTimeout(ctx, k.Limitations.ReadTimeout)
	defer cancel()
	return k.Reader.ReadMessage(ctx)
}

func (k *KafkaClient) Close() {
	k.Reader.Close()
}

func NewKafkaClientFromConfig(config *cviper.CustomViper) KafkaClient {
	kafkaAddress := config.GetStringRequired("QA_HOST") + ":" + config.GetStringRequired("PUBLIC_KAFKA_BROKER_PORT")
	readerConfig := kafka.ReaderConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    config.GetStringRequired("KAFKA_CODE_RUNNER_TOPIC"),
		GroupID:  config.GetStringRequired("KAFKA_CODE_RUNNER_TEST_CONSUMER_GROUP_NAME"),
		MaxBytes: 10e4, // 10MB
	}
	reader := kafka.NewReader(readerConfig)

	return KafkaClient{
		Reader: reader,
		Limitations: KafkaLimitations{
			ReadTimeout: config.GetMillisecondsDurationRequired("KAFKA_READER_READ_TIMEOUT_MILLISECONDS"),
		},
	}
}
