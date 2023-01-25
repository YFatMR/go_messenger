package decorators

import (
	"context"
	"fmt"

	"github.com/YFatMR/go_messenger/core/pkg/loggers"
	"go.uber.org/zap"
)

type DeferFunc = func(error)

func handleLogging(ctx context.Context, logger *loggers.OtelZapLoggerWithTraceID, functionName string,
	structName string,
) DeferFunc {
	logger.InfoContextNoExport(ctx, fmt.Sprintf("`%s` function for `%s` processing...", functionName, structName))
	deferFunc := func(err error) {
		if err != nil {
			logger.ErrorContext(
				ctx, fmt.Sprintf("`%s` function for `%s` failed", functionName, structName), zap.Error(err),
			)
		} else {
			logger.InfoContextNoExport(
				ctx, fmt.Sprintf("`%s` function for `%s` successfully finished...", functionName, structName),
			)
		}
	}
	return deferFunc
}

func LogCallbackError(ctx context.Context, logger *loggers.OtelZapLoggerWithTraceID, functionName string,
	structName string, callback func() error,
) (err error) {
	deferFunc := handleLogging(ctx, logger, "Create", "userRepository")
	defer deferFunc(err)
	return callback()
}

func LogCallbackErrorWithReturnType[T any](ctx context.Context, logger *loggers.OtelZapLoggerWithTraceID,
	functionName string, structName string, callback func() (T, error),
) (_ T, err error) {
	deferFunc := handleLogging(ctx, logger, "Create", "userRepository")
	defer deferFunc(err)
	return callback()
}

func LogCallbackErrorWithTwoReturnType[T any, G any](ctx context.Context, logger *loggers.OtelZapLoggerWithTraceID,
	functionName string, structName string, callback func() (T, G, error),
) (_ T, _ G, err error) {
	deferFunc := handleLogging(ctx, logger, "Create", "userRepository")
	defer deferFunc(err)
	return callback()
}
