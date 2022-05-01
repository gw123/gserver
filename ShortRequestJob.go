package gserver

import (
	"context"
	"fmt"
	"github.com/gw123/gserver/contracts"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"time"
)

type ShortConnJob struct {
	ctx        context.Context
	timeout    time.Duration
	timer      *time.Timer
	cancelFunc context.CancelFunc
	stopFlag   bool
	conn       net.Conn
}

func NewShortConnJob(conn net.Conn, timeout time.Duration, ctx context.Context) *ShortConnJob {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	requestJob := &ShortConnJob{
		timeout:    timeout,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		conn:       conn,
	}
	return requestJob
}

func (job *ShortConnJob) Run(ctx context.Context) error {
	if job.conn == nil {
		return errors.New("conn is nil")
	}
	defer job.conn.Close()
	packer := NewDataPack()
	msg, err := packer.UnPackFromConn(job.conn)
	if err != nil {
		return err
	}
	msg = contracts.NewMsg(1, []byte("server"))
	buf, err := packer.Pack(msg)
	if err != nil {
		fmt.Println("net.Dial 连接错误: ", err)
		return err
	}
	sleep := rand.Int31n(6)
	time.Sleep(time.Duration(sleep) * time.Second)
	job.conn.Write(buf)
	return err
}
