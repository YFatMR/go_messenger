package entity

import (
	"time"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageID struct {
	ID uint64
}

type OffserType string

const (
	BEFORE         = OffserType("before")
	BEFORE_INCLUDE = OffserType("before_include")
	AFTER          = OffserType("after")
	AFTER_INCLUDE  = OffserType("after_include")
)

func OffserTypeFromProtobuf(offset proto.GetDialogMessagesRequest_OffsetType) OffserType {
	if offset == proto.GetDialogMessagesRequest_BEFORE {
		return BEFORE
	} else if offset == proto.GetDialogMessagesRequest_BEFORE_INCLUDE {
		return BEFORE_INCLUDE
	} else if offset == proto.GetDialogMessagesRequest_AFTER {
		return AFTER
	}
	return AFTER_INCLUDE
}

type DialogMessage struct {
	MessageID MessageID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
	Viewed    bool
}

func MessageIDToProtobuf(msg *MessageID) *proto.MessageID {
	return &proto.MessageID{
		ID: msg.ID,
	}
}

func MessageIDFromProtobuf(msg *proto.MessageID) (*MessageID, error) {
	if msg == nil || msg.GetID() == 0 {
		return nil, ErrWrongRequestFormat
	}
	return &MessageID{
		ID: msg.ID,
	}, nil
}

func CopyDialogMessage(msg *DialogMessage) *DialogMessage {
	return &DialogMessage{
		MessageID: msg.MessageID,
		SenderID:  msg.SenderID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
		Viewed:    msg.Viewed,
	}
}

func DialogMessageFromProtobuf(request *proto.CreateDialogMessageRequest) (
	*DialogMessage, error,
) {
	if request == nil || request.GetDialogID().GetID() == 0 || request.GetText() == "" {
		return nil, ErrWrongRequestFormat
	}
	return &DialogMessage{
		Text: request.Text,
	}, nil
}

func DialogMessageToProtobuf(request *DialogMessage, selfID *UserID) *proto.DialogMessage {
	return &proto.DialogMessage{
		MessageID:   MessageIDToProtobuf(&request.MessageID),
		SenderID:    UserIDToProtobuf(&request.SenderID),
		Text:        request.Text,
		CreatedAt:   timestamppb.New(request.CreatedAt),
		SelfMessage: selfID.ID == request.SenderID.ID,
		Viewed:      request.Viewed,
	}
}

func DialogMessagesToProtobuf(messages []*DialogMessage, selfID *UserID) []*proto.DialogMessage {
	result := make([]*proto.DialogMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, DialogMessageToProtobuf(message, selfID))
	}
	return result
}
