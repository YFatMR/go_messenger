package apientity

import (
	"context"

	"github.com/YFatMR/go_messenger/core/pkg/ckafka"
)

type KafkaClient interface {
	WriteNewDialogMessage(ctx context.Context, inMsg *ckafka.DialogMessage) (
		err error,
	)
	Stop()
}
