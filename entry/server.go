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
