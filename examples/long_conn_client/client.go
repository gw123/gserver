package main

import (
	"flag"
	"github.com/gw123/glog"
	"github.com/gw123/gserver"
	"github.com/gw123/gserver/contracts"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var gIsClose = false

func main() {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT)
		sig := <-ch
		switch sig {
		case syscall.SIGINT:
			glog.Info("收到 int信号 程序退出")
			gIsClose = true
		case syscall.SIGTERM:
			glog.Info("收到 term信号 程序强制退出")
			os.Exit(0)
		}
	}()

	glog.Info("client running....")
	workerNum := flag.Int("worker_num", 1, "工作协程数量")
	flag.Parse()

	clientConfig := gserver.LoadClientConfig()
	packer := gserver.NewDataPack()

	glog.Debugf("循环发送消息,workerNum : %d", *workerNum)
	waitGroup := sync.WaitGroup{}
	for i := 0; i < *workerNum; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			client := gserver.NewClient(clientConfig.ServerAddr, clientConfig.Timeout, packer)
			err := client.Connect()
			if err != nil {
				glog.Error(err.Error())
				return
			}
			for !gIsClose {
				msg := contracts.NewMsg(1, []byte("hello world"))
				err = client.Send(msg)
				if err != nil {
					glog.Errorf("%s", err.Error())
					return
				}
				msg, err = client.Read()
				if err != nil {
					glog.Errorf("%s", err.Error())
					return
				}
				// glog.Infof("msgType %v ,msgLen %v, Response:%s", msg.MsgId, msg.Length, string(msg.Body))
			}
		}()
	}
	waitGroup.Wait()
}
