package main

import (
	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/YFatMR/go_messenger/core/pkg/workerpool"
)

func WorkerPoolFromCongig(config *cviper.CustomViper) *workerpool.WorkerPool {
	workersCount := config.GetIntRequired("WORKER_POOL_WORKERS_COUNT")
	return workerpool.New(workersCount)
}
