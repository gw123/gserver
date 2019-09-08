package connpool

import (
	"context"
	"github.com/gw123/glog"
	"github.com/gw123/gserver/contracts"
	"sync"
	"time"
)

type WorkerPool struct {
	status               int
	ctx                  context.Context
	timeout              time.Duration
	Workers              []contracts.Worker
	index                int
	cancelFunc           context.CancelFunc
	poolSize             int
	waitGroup            *sync.WaitGroup
	freePool             chan contracts.Worker
	workerStatusMap      map[int]int
	workerStatusMapMutex sync.RWMutex
}

func NewWorkerPool(ctx context.Context, timeout time.Duration, poolsize int) *WorkerPool {
	ctx, cancelFunc := context.WithCancel(ctx)
	pool := &WorkerPool{
		ctx:             ctx,
		timeout:         timeout,
		Workers:         make([]contracts.Worker, 0),
		cancelFunc:      cancelFunc,
		poolSize:        poolsize,
		waitGroup:       &sync.WaitGroup{},
		freePool:        make(chan contracts.Worker, poolsize),
		workerStatusMap: make(map[int]int),
	}
	pool.init()
	return pool
}

func (pool *WorkerPool) init() {
	for i := 0; i < pool.poolSize; i++ {
		pool.waitGroup.Add(1)
		worker := NewWorker(i, pool.ctx, pool.timeout, pool.waitGroup, pool)
		pool.freePool <- worker
		pool.Workers = append(pool.Workers, worker)
		pool.workerStatusMap[i] = contracts.Worker_Free
	}
}

func (pool *WorkerPool) Push(job contracts.Job) error {
	worker := <-pool.freePool
	worker.Push(job)
	pool.workerStatusMapMutex.Lock()
	pool.workerStatusMap[worker.GetWorkerId()] = contracts.Worker_Running
	pool.workerStatusMapMutex.Unlock()
	return nil
}

func (pool *WorkerPool) pushByForeach(job contracts.Job) error {
	if pool.index >= len(pool.Workers) {
		pool.index = 0
	}

	for i := 0; i < pool.poolSize/2; i++ {
		if pool.Workers[pool.index].IsBusy() {
			pool.index++
			if pool.index >= len(pool.Workers) {
				pool.index = 0
			}
		} else {
			break
		}
	}

	pool.Workers[pool.index].Push(job)
	return nil
}

func (pool *WorkerPool) Run() {
	for i := 0; i < pool.poolSize; i++ {
		go pool.Workers[i].Run()
	}
}

func (pool *WorkerPool) Stop() {
	pool.cancelFunc()
	glog.Debug("wait all worker stop")
	for _, worker := range pool.Workers {
		worker.Stop()
	}
	pool.waitGroup.Wait()
}

/***
  在worker完成一个job后,将其放回free队列等待接收下一个任务
*/
func (pool *WorkerPool) RecycleWorker(worker contracts.Worker) {
	pool.freePool <- worker
	pool.workerStatusMapMutex.Lock()
	pool.workerStatusMap[worker.GetWorkerId()] = contracts.Worker_Free
	pool.workerStatusMapMutex.Unlock()
}

func (pool *WorkerPool) Status() contracts.PoolStatus {
	pool.workerStatusMapMutex.RLock()
	defer pool.workerStatusMapMutex.RUnlock()
	status := contracts.PoolStatus{
		WorkerStausMap: pool.workerStatusMap,
		FreeNum:        len(pool.freePool),
	}

	return status
}
