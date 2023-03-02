package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type KafkaClient interface {
	WriteProgramExecutionMessage(ctx context.Context, programID *entity.ProgramID, userID *entity.UserID) (
		err error,
	)
	Stop()
}
