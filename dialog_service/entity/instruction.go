package entity

import (
	"time"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InstructionOffserType string

const (
	INSTRUCTION_OFFSET_AFTER = InstructionOffserType("after")
)

func InstructionOffserTypeFromProtobuf(offset proto.GetInstructionsByIDRequest_OffsetType) InstructionOffserType {
	// single option at the moment
	return INSTRUCTION_OFFSET_AFTER
}

type InstructionID struct {
	ID uint64
}

func InstructionIDFromProtobuf(instructionID *proto.InstructionID) (*InstructionID, error) {
	if instructionID == nil || instructionID.GetID() == 0 {
		return nil, ErrWrongRequestFormat
	}
	return &InstructionID{
		ID: instructionID.ID,
	}, nil
}

func InstructionIDToProtobuf(instructionID *InstructionID) *proto.InstructionID {
	return &proto.InstructionID{
		ID: instructionID.ID,
	}
}

type Instruction struct {
	InstructionID InstructionID
	CreatedAt     time.Time
	Title         string
	Text          string
}

func InstructionToProtobuf(instruction *Instruction) *proto.Instruction {
	return &proto.Instruction{
		InstructionID: InstructionIDToProtobuf(&instruction.InstructionID),
		CreatedAt:     timestamppb.New(instruction.CreatedAt),
		Title:         instruction.Title,
		Text:          instruction.Text,
	}
}

func InstructionsToProtobuf(instructions []*Instruction) []*proto.Instruction {
	result := make([]*proto.Instruction, 0, len(instructions))
	for _, instruction := range instructions {
		result = append(result, InstructionToProtobuf(instruction))
	}
	return result
}

type CreateInstructionRequest struct {
	InstructionTitle string
	MessagesID       []MessageID
}
