package main

import (
	"github.com/gw123/glog"
	"github.com/gw123/gserver"
)

func RunServer() {
	//LoadServerConfig 可以自己实现,方便扩展或者融合到其他项目中
	config := gserver.LoadServerConfig()
	glog.Infof("%+v", config)
	server := gserver.NewServer(config)
	server.Run()
	glog.Warn("服务正常关闭")
}

func main() {
	RunServer()
}
