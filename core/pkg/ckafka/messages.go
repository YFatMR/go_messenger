package ckafka

import "time"

type ProgramExecutionMessage struct {
	ProgramID string `json:"programId"`
}

type MessageID struct {
	ID uint64 `json:"ID"`
}

type UserID struct {
	ID uint64 `json:"ID"`
}

type DialogID struct {
	ID uint64 `json:"ID"`
}

type DialogMessage struct {
	MessageID MessageID `json:"messageID"`
	SenderID  UserID    `json:"senderID"`
	ReciverID UserID    `json:"reciverID"`
	DialogID  DialogID  `json:"dialogID"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type ViewedMessage struct {
	SenderID         UserID    `json:"senderID"`
	ReciverID        UserID    `json:"reciverID"`
	DialogID         DialogID  `json:"dialogID"`
	MessageID        MessageID `json:"messageID"`
	MessageCreatedAt time.Time `json:"messageCreatedAt"`
}
