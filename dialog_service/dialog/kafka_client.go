package dialog

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaClientSettings struct {
	WriteOperationTimeout time.Duration
}

type KafkaClient struct {
	writer   *kafka.Writer
	settings *KafkaClientSettings
	logger   *czap.Logger
}

func NewKafkaClient(writer *kafka.Writer, settings *KafkaClientSettings, logger *czap.Logger) apientity.KafkaClient {
	return &KafkaClient{
		writer:   writer,
		settings: settings,
		logger:   logger,
	}
}

func (c *KafkaClient) Stop() {
	c.writer.Close()
}

func (c *KafkaClient) WriteNewDialogMessage(ctx context.Context, inMsg *ckafka.DialogMessage) error {
	message, err := json.Marshal(inMsg)
	if err != nil {
		c.logger.ErrorContext(ctx, "Unable to create message", zap.Error(err))
		return ErrMessageCreation
	}

	ctx, cancel := context.WithTimeout(ctx, c.settings.WriteOperationTimeout)
	defer cancel()
	err = c.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(strconv.FormatUint(inMsg.DialogID.ID, 10)),
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
