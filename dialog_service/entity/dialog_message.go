package entity

import (
	"time"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageID struct {
	ID uint64
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

func MessagesIDFromProtobuf(messagesID []*proto.MessageID) ([]*MessageID, error) {
	if messagesID == nil {
		return nil, ErrWrongRequestFormat
	}
	result := make([]*MessageID, 0, len(messagesID))
	for _, messageID := range messagesID {
		msg, err := MessageIDFromProtobuf(messageID)
		if err != nil {
			return nil, err
		}
		result = append(result, msg)
	}
	return result, nil
}

type DialogMessagesOffserType string

const (
	DIALOG_MESSAGE_OFFSET_BEFORE         = DialogMessagesOffserType("before")
	DIALOG_MESSAGE_OFFSET_BEFORE_INCLUDE = DialogMessagesOffserType("before_include")
	DIALOG_MESSAGE_OFFSET_AFTER          = DialogMessagesOffserType("after")
	DIALOG_MESSAGE_OFFSET_AFTER_INCLUDE  = DialogMessagesOffserType("after_include")
)

func DialogMessagesOffserTypeFromProtobuf(offset proto.GetDialogMessagesRequest_OffsetType) DialogMessagesOffserType {
	if offset == proto.GetDialogMessagesRequest_BEFORE {
		return DIALOG_MESSAGE_OFFSET_BEFORE
	} else if offset == proto.GetDialogMessagesRequest_BEFORE_INCLUDE {
		return DIALOG_MESSAGE_OFFSET_BEFORE_INCLUDE
	} else if offset == proto.GetDialogMessagesRequest_AFTER {
		return DIALOG_MESSAGE_OFFSET_AFTER
	}
	return DIALOG_MESSAGE_OFFSET_AFTER_INCLUDE
}

type DialogMessagesType = uint64

const (
	MESSAGE_TYPE_NORMAL = DialogMessagesType(1)
	MESSAGE_TYPE_CODE   = DialogMessagesType(2)
)

func DialogMessagesTypeFromUint64(dialogType uint64) DialogMessagesType {
	switch dialogType {
	case 0:
		return MESSAGE_TYPE_NORMAL
	case 1:
		return MESSAGE_TYPE_CODE
	default:
		return MESSAGE_TYPE_NORMAL
	}
}

func DialogMessagesTypeFromProtobuf(dialogType proto.DialogMessageType) DialogMessagesType {
	switch dialogType {
	case proto.DialogMessageType_NORMAL:
		return MESSAGE_TYPE_NORMAL
	case proto.DialogMessageType_CODE:
		return MESSAGE_TYPE_CODE
	default:
		return MESSAGE_TYPE_NORMAL
	}
}

func DialogMessagesTypeToProtobuf(dialogType DialogMessagesType) proto.DialogMessageType {
	switch dialogType {
	case MESSAGE_TYPE_NORMAL:
		return proto.DialogMessageType_NORMAL
	case MESSAGE_TYPE_CODE:
		return proto.DialogMessageType_CODE
	default:
		return proto.DialogMessageType_NORMAL
	}
}

type CreateDialogMessageRequest struct {
	Text     string
	SenderID UserID
	DialogID DialogID
}

func CreateDialogMessageRequestFromProtobuf(request *proto.CreateDialogMessageRequest) (
	*CreateDialogMessageRequest, error,
) {
	if request == nil || request.GetDialogID().GetID() == 0 || request.GetText() == "" {
		return nil, ErrWrongRequestFormat
	}
	return &CreateDialogMessageRequest{
		Text: request.Text,
		DialogID: DialogID{
			ID: request.DialogID.ID,
		},
	}, nil
}

type CreateDialogMessageWithCodeRequest struct {
	Title    string
	Text     string
	SenderID UserID
	DialogID DialogID
}

func CreateDialogMessageWithCodeRequestFromProtobuf(request *proto.CreateDialogMessageWithCodeRequest) (
	*CreateDialogMessageWithCodeRequest, error,
) {
	if request == nil || request.GetDialogID().GetID() == 0 || request.GetText() == "" || request.GetTitle() == "" {
		return nil, ErrWrongRequestFormat
	}
	return &CreateDialogMessageWithCodeRequest{
		Title: request.Title,
		Text:  request.Text,
		DialogID: DialogID{
			ID: request.DialogID.ID,
		},
	}, nil
}

type DialogMessage struct {
	MessageID MessageID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
	Viewed    bool
	Type      DialogMessagesType
}

func DialogMessageToProtobuf(request *DialogMessage, selfID *UserID) *proto.DialogMessage {
	return &proto.DialogMessage{
		MessageID:   MessageIDToProtobuf(&request.MessageID),
		SenderID:    UserIDToProtobuf(&request.SenderID),
		Text:        request.Text,
		CreatedAt:   timestamppb.New(request.CreatedAt),
		SelfMessage: selfID.ID == request.SenderID.ID,
		Viewed:      request.Viewed,
		Type:        DialogMessagesTypeToProtobuf(request.Type),
	}
}

func DialogMessagesToProtobuf(messages []*DialogMessage, selfID *UserID) []*proto.DialogMessage {
	result := make([]*proto.DialogMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, DialogMessageToProtobuf(message, selfID))
	}
	return result
}
