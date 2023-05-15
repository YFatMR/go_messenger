package sandbox

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type sandboxService struct {
	repository  apientity.SandboxRepository
	kafkaClient apientity.KafkaClient
	logger      *czap.Logger
}

func NewService(repository apientity.SandboxRepository, kafkaClient apientity.KafkaClient, logger *czap.Logger,
) apientity.SandboxService {
	return &sandboxService{
		repository:  repository,
		kafkaClient: kafkaClient,
		logger:      logger,
	}
}

func (s *sandboxService) GetProgramByID(ctx context.Context, programID *entity.ProgramID) (
	*entity.Program, error,
) {
	return s.repository.GetProgramByID(ctx, programID)
}

func (s *sandboxService) CreateProgram(ctx context.Context, programSource *entity.ProgramSource) (
	*entity.ProgramID, error,
) {
	return s.repository.CreateProgram(ctx, programSource)
}

func (s *sandboxService) UpdateProgramSource(ctx context.Context, programID *entity.ProgramID,
	programSource *entity.ProgramSource,
) error {
	return s.repository.UpdateProgramSource(ctx, programID, programSource)
}

func (s *sandboxService) RunProgram(ctx context.Context, programID *entity.ProgramID, userID *entity.UserID) error {
	program, err := s.repository.GetProgramByID(ctx, programID)
	if err != nil {
		return err
	}

	return s.kafkaClient.WriteCodeRunnerMessage(
		ctx, userID, programID, program.Source.SourceCode, program.Source.Language,
	)
}
