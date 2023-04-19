package pgxdb

import (
	"context"
	"fmt"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Connect(ctx context.Context, postgresURL string, logger *czap.Logger) (*pgxpool.Pool, error) {
	for i := 0; i < 10; i++ {
		logger.Debug("Try to connect to databse", zap.Int("try", i+1))
		time.Sleep(2 * time.Second)

		pgConfig, err := pgxpool.ParseConfig(postgresURL)
		if err != nil {
			logger.Error("Failed to parse config database", zap.Error(err))
			continue
		}

		connPool, err := pgxpool.NewWithConfig(ctx, pgConfig)
		if err != nil {
			logger.Error("Failed to create pool to database", zap.Error(err))
			continue
		}

		err = connPool.Ping(ctx)
		if err != nil {
			logger.Error("Failed to ping database", zap.Error(err))
			continue
		}
		return connPool, nil
	}
	return nil, fmt.Errorf("failed to connect to DB")
}
