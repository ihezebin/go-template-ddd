package worker

import (
	"context"
	"sync"

	"github.com/ihezebin/oneness/logger"
)

type workerKeeper struct {
	wg      *sync.WaitGroup
	workers []Worker
}

var keeper *workerKeeper

func init() {
	keeper = &workerKeeper{
		wg:      new(sync.WaitGroup),
		workers: make([]Worker, 0),
	}
}

type Worker interface {
	Name() string
	Run(ctx context.Context) error
	// Cancel 一般不再接受新任务且将未处理完的任务完成
	Cancel()
}

func Run(ctx context.Context) {
	for _, worker := range keeper.workers {
		keeper.wg.Add(1)
		go func(worker Worker) {
			if err := worker.Run(ctx); err != nil {
				logger.WithError(err).Errorf(ctx, "worker [%s] run err", worker.Name())
			}
		}(worker)
	}
}

// Wait 等待 worker 依次平滑退出
func Wait(ctx context.Context) {
	for _, worker := range keeper.workers {
		worker.Cancel()
		keeper.wg.Done()
	}
	logger.Info(ctx, "all workers closed")
	keeper.wg.Wait()
}

func Register(workers ...Worker) {
	keeper.workers = append(keeper.workers, workers...)
}
