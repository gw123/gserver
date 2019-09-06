package connpool

import (
	"context"
	"github.com/gw123/gserver/contracts"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Worker struct {
	stopFlag   bool
	workerId   int
	size       int
	ctx        context.Context
	timeout    time.Duration
	waitTime   time.Duration
	cancelFunc context.CancelFunc
	job        chan contracts.Job
	waitGroup  *sync.WaitGroup
	pool       contracts.WorkerPool
}

func NewWorker(id int,
	parent context.Context,
	timeout time.Duration,
	waitGroup *sync.WaitGroup,
	pool contracts.WorkerPool,
) *Worker {
	ctx, cancelFunc := context.WithCancel(parent)
	jobsize := 10
	worker := &Worker{
		stopFlag:   false,
		workerId:   id,
		ctx:        ctx,
		timeout:    timeout,
		cancelFunc: cancelFunc,
		job:        make(chan contracts.Job, jobsize),
		size:       jobsize,
		waitGroup:  waitGroup,
		pool:       pool,
	}

	return worker
}

func (w *Worker) Push(job contracts.Job) error {
	if w.stopFlag {
		return errors.New("worker stop")
	}
	if job == nil {
		return errors.New("job is nil")
	}
	w.job <- job
	return nil
}

func (w *Worker) Run() {
	defer func() {
		for ; len(w.job) > 0; {
			job := <-w.job
			job.Run()
		}
		if w.waitGroup != nil {
			w.waitGroup.Done()
		}
	}()

	for ; ; {
		if w.stopFlag {
			return
		}

		select {
		case job, ok := <-w.job:
			if !ok {
				return
			}
			job.Run()
			w.pool.RecycleWorker(w)
		default:
			time.Sleep(time.Millisecond * 1)
			//glog.Debug("没有任务 sleep")
		}
	}

}

func (w *Worker) Stop() {
	w.stopFlag = true
}

func (w *Worker) GetTotalJob() int {
	return len(w.job)
}

func (w *Worker) IsBusy() bool {
	return len(w.job) == w.size
}

func (w *Worker) GetWorkerId() int {
	return w.workerId
}
