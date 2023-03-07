package workerpool

import "sync"

type Task struct {
	Execute func()
}

type worker struct {
	taskCh <-chan Task
	done   func()
}

func (w *worker) start() {
	defer w.done()
	for task := range w.taskCh {
		task.Execute()
	}
}

type WorkerPool struct {
	workers []worker
	wg      sync.WaitGroup
	taskCh  chan Task
}

func New(workerCount int) *WorkerPool {
	pool := WorkerPool{
		workers: make([]worker, workerCount),
		taskCh:  make(chan Task),
	}

	pool.wg.Add(workerCount)
	var setupWg sync.WaitGroup
	setupWg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		pool.workers[i] = worker{
			done:   pool.wg.Done,
			taskCh: pool.taskCh,
		}
		go func(i int) {
			setupWg.Done()
			pool.workers[i].start()
		}(i)
	}
	setupWg.Wait()
	return &pool
}

func (w *WorkerPool) AddTask(task Task) error {
	select {
	case w.taskCh <- task:
		return nil
	default:
		return ErrNoFreeWorker
	}
}

func (w *WorkerPool) Stop() {
	close(w.taskCh)
	w.wg.Wait()
}
