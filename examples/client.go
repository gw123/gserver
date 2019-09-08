package main

import (
	"flag"
	"github.com/gw123/glog"
	"github.com/gw123/gserver/client"
	"github.com/gw123/gserver/config"
	"github.com/gw123/gserver/contracts"
	"github.com/gw123/gserver/packdata"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	glog.Info("client runing....")
	workerNum := flag.Int("worker_num", 1, "工作协程数量")
	isLoop := flag.String("loop", "false", "是否循环请求")
	//sleepTime := flag.Int("sleep", 1, "请求等待时间 单位是1s 默认是1s")
	flag.Parse()

	clientConfig := config.LoadClientConfig()
	glog.Dump(clientConfig)
	signer := packdata.NewSignerHashSha1([]byte(clientConfig.Key))
	packer := packdata.NewDataPackV1(signer)

	if *isLoop == "true" {
		glog.Debug("循环发送消息,workerNum : %d", *workerNum)
		waitGroup := sync.WaitGroup{}
		for i := 0; i < *workerNum; i++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for ; !gIsClose; {
					client := client.NewClient(clientConfig.ServerAddr, clientConfig.Timeout, packer)
					err := client.Connect()
					if err != nil {
						glog.Error(err.Error())
						return
					}

					msg := contracts.NewMsg(1, []byte("hello world"))
					err = client.Send(msg)
					if err != nil {
						glog.Error("%s", err.Error())
						return
					}
					msg, err = client.Read()
					if err != nil {
						glog.Error("%s", err.Error())
						return
					}
					glog.Debug("Response:%s", string(msg.Body))
					time.Sleep(time.Millisecond)
				}
			}()
		}
		waitGroup.Wait()
	} else {
		client := client.NewClient(clientConfig.ServerAddr, clientConfig.Timeout, packer)
		err := client.Connect()
		if err != nil {
			glog.Error(err.Error())
			return
		}
		glog.Debug("发送一条消息,workerNum : %d", *workerNum)
		msg := contracts.NewMsg(1, []byte("hello world"))
		err = client.Send(msg)
		if err != nil {
			glog.Error("%s", err.Error())
			return
		}
	}
}
