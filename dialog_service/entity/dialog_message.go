package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type DialogMessage struct {
	SenderID  UserID
	DialogID  DialogID
	Text      string
	CreatedAt uint64
}

// func DialogMessageFromProtobuf(request *proto.CreateDialogMessageRequest) (
// 	*DialogMessage, error,
// ) {
// 	if request == nil || request.GetDialogID().GetID() == 0 || request.GetText() == "" ||
// 		request.GetSenderID().GetID() == "" {
// 		return nil, ErrWrongRequestFormat
// 	}

// 	senderID, err := UserIDFromProtobuf(request.SenderID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dialogID, err := DialogIDFromProtobuf(request.DialogID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DialogMessage{
// 		SenderID: *senderID,
// 		DialogID: *dialogID,
// 		Text:     request.Text,
// 	}, nil
// }

func DialogMessageToProtobuf(request *DialogMessage) *proto.DialogMessage {
	return &proto.DialogMessage{
		SenderID:  UserIDToProtobuf(&request.SenderID),
		DialogID:  DialogIDToProtobuf(&request.DialogID),
		Text:      request.Text,
		CreatedAt: request.CreatedAt,
	}
}

func DialogMessagesToProtobuf(messages []*DialogMessage) []*proto.DialogMessage {
	result := make([]*proto.DialogMessage, 0, len(messages))
	for _, message := range messages {
		result = append(result, DialogMessageToProtobuf(message))
	}
	return result
}
