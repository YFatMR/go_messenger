package entity

import (
	"fmt"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type BotMessage struct {
	Text string
}

func GPTMessageFromProtobuf(inMessage *proto.BotMessage) (
	*BotMessage, error,
) {
	if inMessage == nil || inMessage.GetText() == "" {
		return nil, fmt.Errorf("expected not nil message content")
	}
	return &BotMessage{
		Text: inMessage.Text,
	}, nil
}

func GPTMessageToProtobuf(inMessage *BotMessage) *proto.BotMessage {
	return &proto.BotMessage{
		Text: inMessage.Text,
	}
}
