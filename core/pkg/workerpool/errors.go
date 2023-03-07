package workerpool

import "errors"

var ErrNoFreeWorker = errors.New("no free worker")
