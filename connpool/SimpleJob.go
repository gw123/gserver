package connpool

import (
	"context"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type SimpleJob struct {
	ctx        context.Context
	timeout    time.Duration
	timer      *time.Timer
	cancelFunc context.CancelFunc
	stopFlag   bool
	data       interface{}
}

func NewSimpleJob(data interface{}, timeout time.Duration, ctx context.Context) *SimpleJob {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	requestJob := &SimpleJob{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		data:       data,
	}
	return requestJob
}

func (job *SimpleJob) Run() (interface{}, error) {
	sleepTime := rand.Int31n(5)
	var i int32 = 0
	for ; i < sleepTime; i++ {
		if job.stopFlag == true {
			return nil, errors.New("break by job stop")
			break
		}

	}
	return job.data, nil
}

func (job *SimpleJob) Stop() {
	job.stopFlag = true
	return
}

func (job *SimpleJob) GetJobType() string {
	job.stopFlag = true
	return "simpleJob"
}
