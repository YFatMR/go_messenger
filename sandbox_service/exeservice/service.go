package exeservice

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/ulo"
	"github.com/YFatMR/go_messenger/sandbox_service/entities/program"
	"github.com/YFatMR/go_messenger/sandbox_service/entities/programresult"
	"github.com/YFatMR/go_messenger/sandbox_service/sandbox"
)

type ProgramExecutionService interface {
	Execute(ctx context.Context, program *program.Entity) (
		programResult *programresult.Entity, logStash ulo.LogStash, err error,
	)
}

type programExecutionService struct {
	programmExecutionTimeout time.Duration
	sandboxClient            sandbox.Client
}

func New(sandboxClient sandbox.Client, programmExecutionTimeout time.Duration) ProgramExecutionService {
	return &programExecutionService{
		programmExecutionTimeout: programmExecutionTimeout,
		sandboxClient:            sandboxClient,
	}
}

func (s *programExecutionService) Execute(ctx context.Context, program *program.Entity) (
	*programresult.Entity, ulo.LogStash, error,
) {
	ctx, cancel := context.WithTimeout(ctx, s.programmExecutionTimeout)
	defer cancel()
	stdout, stderr, err := s.sandboxClient.ExecuteGoCode(ctx, program.GetSourceCode(), "user_id_1")
	if err != nil {
		return nil, ulo.FromError(err), err
	}
	return programresult.New(stdout.String(), stderr.String()), nil, nil
}
