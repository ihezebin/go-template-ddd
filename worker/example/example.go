package example

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ihezebin/go-template-ddd/worker"
	"github.com/pkg/errors"
)

type exampleWorker struct {
	wg         *sync.WaitGroup
	processNum int
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (worker *exampleWorker) Name() string {
	return "example"
}

func (worker *exampleWorker) Run(ctx context.Context) error {
	if worker == nil {
		return errors.New("worker is nil")
	}
	ctx, worker.cancelFunc = context.WithCancel(ctx)
	worker.ctx = ctx

	for i := 0; i < worker.processNum; i++ {
		worker.wg.Add(1)
		go func() {
			defer worker.wg.Done()
			worker.handle()
		}()
	}

	return nil
}

func (worker *exampleWorker) Cancel() {
	worker.cancelFunc()
	worker.wg.Wait()
}

func (worker *exampleWorker) handle() {
	for {
		select {
		case <-worker.ctx.Done():
			return
		default:
			// do something here
			fmt.Printf("%s do something here, %s\n", worker.Name(), time.Now().Format(time.DateTime))
			time.Sleep(time.Hour)
		}
	}
}

var _ worker.Worker = (*exampleWorker)(nil)

func NewExampleWorker() worker.Worker {
	return &exampleWorker{
		wg:         new(sync.WaitGroup),
		processNum: 2,
	}
}
