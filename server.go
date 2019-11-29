package gserver

import (
	"context"
	"github.com/gw123/glog"
	"github.com/gw123/gserver/contracts"
	"github.com/gw123/gworker"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	config     contracts.ServerConfig
	pool       gworker.WorkerPool
	conn       net.Conn
	count      int
	totalCount int
	countMutex sync.Mutex
	isClose    bool
	ctx        context.Context
}

func NewServer(config contracts.ServerConfig) *Server {
	ctx := context.Background()
	workerPool := gworker.NewWorkerPool(ctx, time.Second*2, 100, func(err error, job gworker.Job) {
		log.Print(err)
	})
	return &Server{
		config: config,
		pool:   workerPool,
		ctx:    ctx,
	}
}

func (s *Server) server() {
	addr := s.config.BindAddr
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Error("%s", err.Error())
		return
	}
	glog.Info("server listen at %s", addr)
	defer listen.Close()
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				if !s.isClose {
					glog.Error(err.Error())
				}
				break
			}
			if s.isClose {
				glog.Warn("连接在手动关闭信号后到来")
				conn.Close()
				break
			}

			s.countMutex.Lock()
			s.count++
			s.countMutex.Unlock()
			job := NewRequestJob(conn, time.Second*2, s.ctx)
			s.pool.Push(job)
		}
	}()

	timer := time.NewTicker(time.Second * 5)
	for {
		if s.isClose {
			time.Sleep(time.Second * 2)
			break
		}
		select {
		case <-timer.C:
			glog.Info("5s内一共收到%d个请求", s.count)
			s.countMutex.Lock()
			s.count = 0
			s.countMutex.Unlock()
			status := s.Status()
			glog.Info("freeNum : %d", status)
		}
	}
}

func (s *Server) Run() {
	go s.handleSignal()
	s.pool.Run()
	s.server()
}

func (s *Server) Stop() {
	s.isClose = true
	s.pool.Stop()
}

func (s *Server) Status() uint {
	return s.pool.Status()
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
