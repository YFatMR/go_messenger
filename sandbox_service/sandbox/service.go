package sandbox

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ctrace"
	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/workerpool"
	"github.com/YFatMR/go_messenger/sandbox_service/apientity"
	"github.com/YFatMR/go_messenger/sandbox_service/entity"
	"go.uber.org/zap"
)

type sandboxService struct {
	repository  apientity.SandboxRepository
	codeRunner  apientity.CodeRunner
	workerPool  *workerpool.WorkerPool
	kafkaClient apientity.KafkaClient
	logger      *czap.Logger
}

func NewService(repository apientity.SandboxRepository, codeRunner apientity.CodeRunner,
	workerPool *workerpool.WorkerPool, kafkaClient apientity.KafkaClient, logger *czap.Logger,
) apientity.SandboxService {
	return &sandboxService{
		repository:  repository,
		codeRunner:  codeRunner,
		workerPool:  workerPool,
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
	traceID := ctrace.TraceIDFromContext(ctx)
	program, err := s.repository.GetProgramByID(ctx, programID)
	if err != nil {
		return err
	}

	return s.workerPool.AddTask(workerpool.Task{
		Execute: func() {
			ctx := context.Background()
			output, err := s.codeRunner.RunGoCode(ctx, program.Source.SourceCode, userID.ID)
			s.logger.Debug("Run program", zap.String("program ID", programID.ID))
			if err != nil {
				// TODO: mutex for logger
				s.logger.Error(
					"Can not run program", zap.Error(err), zap.String("program ID", programID.ID),
					zap.String("user ID", userID.ID), zap.String(czap.TraceIDKey, traceID),
				)
				return
			}

			err = s.repository.UpdateCodeRunnerOutput(ctx, programID, output)
			if err != nil {
				s.logger.Error("Can not update program ren result", zap.String(czap.TraceIDKey, traceID))
				return
			}

			s.kafkaClient.WriteProgramExecutionMessage(ctx, programID, userID)
		},
	})
}

// TODO: implement
// LintProgram(ctx context.Context, programID *entity.ProgramID) (
// 	err error,
// )
