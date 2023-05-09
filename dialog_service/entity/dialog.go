package entity

import (
	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
)

type DialogID struct {
	ID uint64
}

type DialogMember struct {
	UserID             UserID
	UnredMessagesCount uint64
}

type Dialog struct {
	DialogID            DialogID
	Name                string
	UnreadMessagesCount uint64
	LastMessage         DialogMessage
	LastReadMessage     DialogMessage
}

func DialogIDFromProtobuf(dialogID *proto.DialogID) (*DialogID, error) {
	if dialogID == nil || dialogID.GetID() == 0 {
		return nil, ErrWrongRequestFormat
	}
	return &DialogID{
		ID: dialogID.GetID(),
	}, nil
}

func DialogIDToProtobuf(dialogID *DialogID) *proto.DialogID {
	return &proto.DialogID{
		ID: dialogID.ID,
	}
}

func DialogToProtobuf(dialog *Dialog, selfID *UserID) *proto.Dialog {
	return &proto.Dialog{
		DialogID:            DialogIDToProtobuf(&dialog.DialogID),
		Name:                dialog.Name,
		UnreadMessagesCount: dialog.UnreadMessagesCount,
		LastMessage:         DialogMessageToProtobuf(&dialog.LastMessage, selfID),
		LastReadMessage:     DialogMessageToProtobuf(&dialog.LastReadMessage, selfID),
	}
}

func DialogsToProtobuf(dialogs []*Dialog, selfID *UserID) []*proto.Dialog {
	result := make([]*proto.Dialog, 0, len(dialogs))
	for _, dialog := range dialogs {
		result = append(result, DialogToProtobuf(dialog, selfID))
	}
	return result
}
