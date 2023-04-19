package sandbox

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/decorator"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaClient struct {
	writer                *kafka.Writer
	logger                *czap.Logger
	writeOperationTimeout time.Duration
}

func KafkaClientFromConfig(config *cviper.CustomViper, logger *czap.Logger) apientity.KafkaClient {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.GetStringRequired("KAFKA_BROKER_ADDRESS")),
		Topic:        config.GetStringRequired("KAFKA_CODE_RUNNER_TOPIC"),
		Balancer:     &kafka.LeastBytes{},
		Compression:  kafka.Snappy,
		WriteTimeout: config.GetMillisecondsDurationRequired("KAFKA_WRITER_WRITE_TIMEOUT_MILLISECONDS"),
		ReadTimeout:  config.GetMillisecondsDurationRequired("KAFKA_WRITER_READ_TIMEOUT_MILLISECONDS"),
	}
	base := &KafkaClient{
		writer: writer,
		logger: logger,
		// TODO: check that it's true
		writeOperationTimeout: config.GetMillisecondsDurationRequired("KAFKA_WRITER_WRITE_TIMEOUT_MILLISECONDS"),
	}
	return decorator.NewLoggingKafkaClientDecorator(base, logger)
}

func (c *KafkaClient) Stop() {
	c.writer.Close()
}

func (c *KafkaClient) WriteProgramExecutionMessage(ctx context.Context, programID *entity.ProgramID,
	userID *entity.UserID,
) error {
	message, err := json.Marshal(ckafka.ProgramExecutionMessage{
		ProgramID: programID.ID,
	})
	if err != nil {
		c.logger.ErrorContext(ctx, "Unable to create message", zap.Error(err))
		return ErrMessageCreation
	}

	ctx, cancel := context.WithTimeout(ctx, c.writeOperationTimeout)
	defer cancel()
	err = c.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(strconv.FormatUint(userID.ID, 10)),
			Value: message,
			Time:  time.Now(),
		},
	)
	if err != nil {
		c.logger.ErrorContext(ctx, "Unable to write message", zap.Error(err))
		return ErrMessageWriting
	}
	return nil
}
