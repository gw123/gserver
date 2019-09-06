package connpool

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPoll(t *testing.T) {
	ctx := context.Background()
	pool := NewWorkerPool(ctx, time.Second*2, 100)
	pool.Run()
	mutex := sync.Mutex{}
	count := 0
	go func() {
		for count = 0; ; {
			mutex.Lock()
			count++
			mutex.Unlock()
			pool.Push(nil)
		}
	}()

	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-timer.C:
			t.Log(count / 5)
			mutex.Lock()
			count = 0
			mutex.Unlock()

		}
	}
}

