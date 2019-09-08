#!/usr/bin/env bash
srcDir=${GOPATH}"/src/github.com/gw123/gserver"
dstDir=${GOPATH}"/bin/gserver"

if [ ! -d "${dstDir}" ]; then
    mkdir -p ${dstDir}
fi

term() {
    echo "term"
    ps -aux|grep pos-tcpserver|grep -v grep|awk '{print $2}'|xargs kill -TERM
}

server() {
    buildStr="go build -o  ${dstDir}/server  ${srcDir}/server.go"
    echo $buildStr
    $buildStr
    if [ $? -eq 0 ]; then
        cd $dstDir && ./server
    else
       echo "编译失败"
    fi
}

client() {
    buildStr="go build -o  ${dstDir}/client  ${srcDir}/client.go"
    echo $buildStr
    $buildStr
    if [ $? -eq 0 ]; then
        cd $dstDir && ./client
    else
       echo "编译失败"
    fi
}


runDocker(){
    docker run -it --rm -v $PWD/examples:/etc/gserver\
    -p 8080:8080\
    ccr.ccs.tencentyun.com/g-docker/gserver:latest
}


case $1 in
 "term")
 term
  ;;
  "server")
  server
  ;;
  "client")
  client
  ;;
  "runDocker")
  runDocker
  ;;
  *)
   echo "$0 [server|client|runDocker]"
  ;;
esac
