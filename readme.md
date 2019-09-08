#基于tcp协议的服务和客户端  
- 使用tlv协议头只需要4或者8字节
- 实现协程池避免协程创建和释放的开销,可以稳定支持几w的并非连接
- 实现平滑关闭在接受到关闭信号后等待所有的woker都执行完毕后才会关闭
- 实现拆包和解包和包内容安全校验 ,支持hash_mac_sha1,hash_mac_sha256,hash_mac_sha512签名验证
- 采用worker do jobs 模式 ,方便扩展对conn 不同的处理方式 
- 实现worker 负载监控,当worker 消耗超过 指定值触发警告(可以配置向一个接口发一个请求, 发邮件) todo
- todo 监控每个job 运行时长
- todo 将日志和监控发送到elk中
- todo 使用etcd 动态更新配置

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
	"github.com/gw123/gserver/config"
	"github.com/gw123/gserver/server"
)

func RunServer() {
	config := config.LoadServerConfig()
	glog.Dump(config)
	server := server.NewServer(config)
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
	flag.Parse()

	clientConfig := config.LoadClientConfig()
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
```
## go run entry/client.go 

## 其他可以参考 examples