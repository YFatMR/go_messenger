package ckafka

type ProgramExecutionMessage struct {
	ProgramID string `json:"programId"`
	UserID    string `json:"userId"`
}
