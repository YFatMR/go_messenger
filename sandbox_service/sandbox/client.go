package sandbox

import (
	"bytes"
	"context"
)

type Client interface {
	ExecuteGoCode(ctx context.Context, sourceCode string, userID string) (
		stdout *bytes.Buffer, stderr *bytes.Buffer, err error,
	)
}
