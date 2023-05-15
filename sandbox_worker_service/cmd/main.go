package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := cviper.New()
	config.AutomaticEnv()

	executionCommand := config.GetStringRequired("EXECUTION_COMMAND")
	executionCommandArgs := ExecutionCommandArgsFromConfig(config)

	// Init logger
	logger, err := czap.FromConfig(config)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic!", zap.Any("msg", r))
		}
		panic("Panic")
	}()

	// KAFKA READER
	logger.DebugContext(ctx, "Startin read kafka messages")

	kafkaReader := NewKafkaReaderFromConfig(config)

	logger.DebugContext(ctx, "Readed kafka message")

	kafkaInMessage, err := kafkaReader.ReadMessage(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Can not read messages", zap.Error(err))
		return
	}
	parsedInMessage := ckafka.CodeRunnerMessage{}
	if err = json.Unmarshal(kafkaInMessage.Value, &parsedInMessage); err != nil {
		logger.ErrorContext(ctx, "Can't parse kafka message", zap.Error(err))
		return
	}
	logger.DebugContext(ctx, "Successfully parsed messages")

	// RUN PROGRAM
	logger.DebugContext(ctx, "Starting run program")

	executionCommandArgs = append(executionCommandArgs, parsedInMessage.SourceCode)
	cmd := exec.Command(executionCommand, executionCommandArgs...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if err != nil {
		logger.ErrorContext(ctx, "Failed to run command", zap.Error(err))
		return
	}

	logger.DebugContext(ctx, "Program finished")

	// WRITE PROGRAM RESILT

	logger.DebugContext(ctx, "Starting write result to kafka")

	kafkaWriter := NewKafkaWriterFromConfig(config)
	outKafkaMessage, err := json.Marshal(ckafka.CodeRunnerResultMessage{
		ProgramID: parsedInMessage.ProgramID,
		Stdout:    stdoutBuf.String(),
		Stderr:    stderrBuf.String(),
	})
	if err != nil {
		logger.ErrorContext(ctx, "Unable to create out message", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(ctx, kafkaWriter.WriteTimeout)
	defer cancel()
	err = kafkaWriter.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(strconv.FormatInt(int64(parsedInMessage.SenderID.ID), 10)),
			Value: outKafkaMessage,
			Time:  time.Now(),
		},
	)
	if err != nil {
		logger.ErrorContext(ctx, "Unable to write message", zap.Error(err))
		return
	}
	logger.DebugContext(ctx, "Successfully finished")
}
