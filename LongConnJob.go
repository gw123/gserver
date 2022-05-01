package gserver

import (
	"context"
	"github.com/gw123/glog"
	"github.com/gw123/gserver/contracts"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net"
	"time"
)

type LongConnJob struct {
	timeout    time.Duration
	timer      *time.Timer
	cancelFunc context.CancelFunc
	stopFlag   bool
	conn       net.Conn
}

func NewLongConnJob(conn net.Conn, timeout time.Duration) *LongConnJob {
	requestJob := &LongConnJob{
		timeout: timeout,
		conn:    conn,
	}
	return requestJob
}

func (job *LongConnJob) Run(ctx context.Context) error {
	debug := viper.GetBool("debug")

	if job.conn == nil {
		return errors.New("conn is nil")
	}

	var flag = true
	defer func() {
		flag = false
	}()

	packer := NewDataPack()
	var count = 0

	go func() {
		for flag {
			glog.Infof("count %d", count)
			time.Sleep(time.Second)

			select {
			case <-ctx.Done():
				flag = false
				// todo signal client server has close
				job.conn.Close()
				break
			default:
			}
		}
	}()

	for flag {
		count++
		recvMsg, err := packer.UnPackFromConn(job.conn)
		if err != nil {
			glog.WithErr(err).Errorf("packer unPackFromConn err")
			return err
		}

		if debug {
			glog.Infof("sendMsg msgId", recvMsg.MsgId, recvMsg.Length)
		}

		sendMsg := contracts.NewMsg(1, []byte("server"))
		buf, err := packer.Pack(sendMsg)
		if err != nil {
			glog.WithErr(err).Errorf("packer pack err")
			continue
		}

		job.conn.Write(buf)
	}

	return nil
}
