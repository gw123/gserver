package contracts

import "net"

type Worker interface {
	IsBusy() bool
	Push(job Job) error
	Run()
	Stop()
	GetTotalJob() int
	GetWorkerId() int
}

type WorkerPool interface {
	Push(job Job) error
	Run()
	Stop()
	RecycleWorker(worker Worker)
	Status() PoolStatus
}

type RemoteClient interface {
	GetConn() net.Conn
	Stop()
	HandleRequest() error
}
