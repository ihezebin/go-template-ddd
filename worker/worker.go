package worker

import (
	"context"
	"sync"

	"github.com/ihezebin/soup/logger"
)

type worKeeper struct {
	wg      *sync.WaitGroup
	workers []Worker
}

func (k *worKeeper) Name() string {
	return "worker_keeper"
}

func (k *worKeeper) Run(ctx context.Context) error {
	for _, worker := range k.workers {
		k.wg.Add(1)
		go func(worker Worker) {
			if err := worker.Run(ctx); err != nil {
				logger.WithError(err).Errorf(ctx, "worker [%s] run err", worker.Name())
			}
		}(worker)
	}
	k.wg.Wait()
	return nil
}

func (k *worKeeper) Close(ctx context.Context) error {
	for _, worker := range k.workers {
		worker.Cancel()
		k.wg.Done()
	}
	logger.Info(ctx, "all workers closed")
	return nil
}

type Worker interface {
	Name() string
	Run(ctx context.Context) error
	// Cancel 一般不再接受新任务且将未处理完的任务完成
	Cancel()
}

func NewWorKeeper(workers ...Worker) *worKeeper {
	return &worKeeper{
		wg:      new(sync.WaitGroup),
		workers: workers,
	}
}
