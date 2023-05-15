package ckafka

import "time"

type UserID struct {
	ID uint64 `json:"ID"`
}

type CodeRunnerMessage struct {
	SenderID   UserID `json:"senderID"`
	ProgramID  string `json:"programID"`
	SourceCode string `json:"sourceCode"`
	Language   string `json:"languageWithVersion"`
}

type CodeRunnerResultMessage struct {
	ProgramID string `json:"programID"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
}

type MessageID struct {
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
	Type      uint64    `json:"type"`
}

type ViewedMessage struct {
	SenderID         UserID    `json:"senderID"`
	ReciverID        UserID    `json:"reciverID"`
	DialogID         DialogID  `json:"dialogID"`
	MessageID        MessageID `json:"messageID"`
	MessageCreatedAt time.Time `json:"messageCreatedAt"`
}
