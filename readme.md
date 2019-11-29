#基于tcp协议的服务和客户端  
- 使用tlv协议头只需要8字节
- 实现协程池避免协程创建和释放的开销,可以稳定支持几w的并发连接
- 实现平滑关闭在接受到关闭信号后等待所有的woker都执行完毕后才会关闭
- 采用worker do jobs 模式 ,方便扩展对conn不同的处理方式 


## 在实际使用过程中发现有些客户端实现存在一些问题 dataLen 传过来是一个错误的值 ,这就导致了服务端要创建的buffer过大
出现out of memery错误导致服务异常关闭, 所以为了避免这个问题, 这里限制了每个数据包的长度不大于1M ,超过1M的数据包将会被丢弃
并且返回: "报文长度异常的错误"

# Server 配置 config.server.yaml
```
bind_addr: "127.0.0.1:8080"
key: "123456"
timeout: 15
pool_size: 1000
crypto_type: "hash_sha1"

```

# 代码实现
```
package main

import (
	"github.com/gw123/glog"
	"github.com/gw123/gserver"
)

func RunServer() {
	//LoadServerConfig 可以自己实现,方便扩展或者融合到其他项目中
	config := gserver.LoadServerConfig()
	glog.Dump(config)
	server := gserver.NewServer(config)
	server.Run()
	glog.Warn("服务正常关闭")
}

func main() {
	RunServer()
}

```
## go run  entry/server.go 

#Client 配置 config.client.yaml
```
server_addr: "127.0.0.1:8080"
client_num: 1
key: "123456"
timeout: 15
pool_size: 1000
crypto_type: "hash_sha1"
```

#代码实现
```
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

	clientConfig := gserver.LoadClientConfig()
	glog.Dump(clientConfig)
	packer := gserver.NewDataPack()

	if *isLoop == "true" {
		glog.Debug("循环发送消息,workerNum : %d", *workerNum)
		waitGroup := sync.WaitGroup{}
		for i := 0; i < *workerNum; i++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for ; !gIsClose; {
					client := gserver.NewClient(clientConfig.ServerAddr, clientConfig.Timeout, packer)
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
		client := gserver.NewClient(clientConfig.ServerAddr, clientConfig.Timeout, packer)
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

```
## go run entry/client.go 

## 其他可以参考 examples