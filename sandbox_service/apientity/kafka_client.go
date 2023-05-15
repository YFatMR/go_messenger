package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type KafkaClient interface {
	WriteCodeRunnerMessage(ctx context.Context, userID *entity.UserID, programID *entity.ProgramID,
		sourceCode string, language entity.Languages,
	) (
		err error,
	)
	Stop()
}
