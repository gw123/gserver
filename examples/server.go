package main

import (
	"github.com/gw123/glog"
	"github.com/gw123/gserver/config"
	"github.com/gw123/gserver/server"
	"log"
	"os"
	"runtime/pprof"
)

func RunServer() {
	config := config.LoadServerConfig()
	glog.Dump(config)
	server := server.NewServer(config)
	server.Run()
	glog.Warn("服务正常关闭")
}

func main() {
	//go func() {
	//	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
	//}()
	cpuf, err := os.Create("cpu_profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuf)
	defer pprof.StopCPUProfile()

	//memf, err := os.Create("mem_profile")
	//if err != nil {
	//	log.Fatal("could not create memory profile: ", err)
	//}
	//if err := pprof.WriteHeapProfile(memf); err != nil {
	//	log.Fatal("could not write memory profile: ", err)
	//}
	//defer  memf.Close()

	RunServer()
}
