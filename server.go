package gserver

import (
	"context"
	"github.com/gw123/glog"
	"github.com/gw123/gserver/contracts"
	"github.com/gw123/gworker"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Server struct {
	config     contracts.ServerConfig
	pool       gworker.WorkerPool
	conn       net.Conn
	count      int64
	totalCount int
	countMutex sync.Mutex
	isClose    bool
	ctx        context.Context
	cancel     context.CancelFunc
	listener   net.Listener
}

func NewServer(config contracts.ServerConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	workerPool := gworker.NewWorkerPool(ctx, time.Second*2, 100)

	return &Server{
		config: config,
		pool:   workerPool,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) server() {
	addr := s.config.BindAddr
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Errorf("%s", err.Error())
		return
	}
	s.listener = listen
	glog.Infof("server listen at %s", addr)

	for !s.isClose {
		conn, err := listen.Accept()
		if err != nil {
			glog.WithErr(err).Errorf("accept tcp conn err")
			break
		}

		atomic.AddInt64(&s.count, 1)
		job := NewLongConnJob(conn, time.Second*2)
		s.pool.Push(job)
	}
}

func (s *Server) Run() {
	go s.handleSignal()
	go s.Status()
	s.pool.Run()
	s.server()
}

func (s *Server) Stop() {
	s.isClose = true
	s.pool.Stop()
	s.listener.Close()
}

func (s *Server) Status() {
	timer := time.NewTicker(time.Second * 10)
	for {
		if s.isClose {
			time.Sleep(time.Second * 2)
			break
		}
		select {
		case <-timer.C:
			atomic.StoreInt64(&s.count, 0)
			status := s.pool.Status()
			glog.Infof("freeNum : %d", status)
		}
	}
	return
}

func (s *Server) GetConfig() contracts.ServerConfig {
	return s.config
}

func (s *Server) handleSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	sig := <-ch
	switch sig {
	case syscall.SIGINT:
		glog.Info("收到 int信号 程序退出")
		s.isClose = true
		s.Stop()
	case syscall.SIGTERM:
		glog.Info("收到 term信号 程序强制退出")
		os.Exit(0)
	}
}
