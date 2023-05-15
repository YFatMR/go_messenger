package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
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

type KafkaSettings struct {
	TopicMap map[entity.Languages]string
}

type KafkaClient struct {
	writer                *kafka.Writer
	settings              *KafkaSettings
	logger                *czap.Logger
	writeOperationTimeout time.Duration
}

func KafkaClientFromConfig(config *cviper.CustomViper, logger *czap.Logger,
) apientity.KafkaClient {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.GetStringRequired("KAFKA_BROKER_ADDRESS")),
		Balancer:     &kafka.LeastBytes{},
		Compression:  kafka.Snappy,
		WriteTimeout: config.GetMillisecondsDurationRequired("KAFKA_WRITER_WRITE_TIMEOUT_MILLISECONDS"),
		ReadTimeout:  config.GetMillisecondsDurationRequired("KAFKA_WRITER_READ_TIMEOUT_MILLISECONDS"),
	}
	topicMap := make(map[entity.Languages]string)
	topicMap[entity.PYTHON_V_3_9] = config.GetStringRequired("KAFKA_PYTHON_3_9_CODE_RUNNER_TOPIC")

	base := &KafkaClient{
		writer: writer,
		logger: logger,
		settings: &KafkaSettings{
			TopicMap: topicMap,
		},
		// TODO: check that it's true
		writeOperationTimeout: config.GetMillisecondsDurationRequired("KAFKA_WRITER_WRITE_TIMEOUT_MILLISECONDS"),
	}
	return decorator.NewLoggingKafkaClientDecorator(base, logger)
}

func (c *KafkaClient) Stop() {
	c.writer.Close()
}

func (c *KafkaClient) WriteCodeRunnerMessage(ctx context.Context, userID *entity.UserID, programID *entity.ProgramID,
	sourceCode string, language entity.Languages,
) error {

	message, err := json.Marshal(ckafka.CodeRunnerMessage{
		ProgramID:  programID.ID,
		SourceCode: sourceCode,
		Language:   string(language),
		SenderID: ckafka.UserID{
			ID: userID.ID,
		},
	})
	if err != nil {
		c.logger.ErrorContext(ctx, "Unable to create message", zap.Error(err))
		return ErrMessageCreation
	}

	topic, ok := c.settings.TopicMap[language]
	if !ok {
		c.logger.ErrorContext(ctx, "Unsupported language", zap.String("lang", string(language)))
		return fmt.Errorf("unsupported language")
	}

	ctx, cancel := context.WithTimeout(ctx, c.writeOperationTimeout)
	defer cancel()
	err = c.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(strconv.FormatUint(userID.ID, 10)),
			Value: message,
			Topic: topic,
			Time:  time.Now(),
		},
	)
	if err != nil {
		c.logger.ErrorContext(ctx, "Unable to write message", zap.Error(err))
		return ErrMessageWriting
	}
	return nil
}
