package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"go-im/common/socketio"
	"go-im/connect/internal/config"
	"go-im/connect/internal/logic"
	"go-im/connect/internal/server"
	"go-im/connect/internal/svc"

	zeroservice "github.com/zeromicro/go-zero/core/service"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/connect.yaml", "the config file")

func main() {
	flag.Parse()

	var err error
	var c config.Config
	conf.MustLoad(*configFile, &c)
	srvCtx := svc.NewServiceContext(c)

	logx.DisableStat()

	//tcpServer := server.NewTCPServer(srvCtx)
	wsServer := server.NewWSServer(srvCtx)

	//tcpServer.Server, err = socket.NewServe(c.Name, c.TCPListenOn, c.SendChanSize)
	//if err != nil {
	//	panic(err)
	//}

	wsServer.Server, err = socketio.NewServe(c.Name, c.WSListenOn, c.SendChanSize)
	if err != nil {
		panic(err)
	}
	
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.HandleRequest(w,r)
	})


	go wsServer.Start()
	//go tcpServer.HandleRequest()
	go wsServer.KqHeart()

	fmt.Printf("Starting tcp server at %s, ws server at: %s...\n", c.TCPListenOn, c.WSListenOn)

	serviceGroup := zeroservice.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range logic.Consumers(context.Background(), srvCtx, wsServer.Server) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
