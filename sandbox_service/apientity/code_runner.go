package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/sandbox_service/entity"
)

type CodeRunner interface {
	RunGoCode(ctx context.Context, sourceCode string, userID string) (
		output *entity.ProgramOutput, err error,
	)
	Stop()
	// LintCode(ctx context.Context, sourceCode string, userID string) (
	// 	stdout []byte, stderr []byte, err error,
	// )
}
