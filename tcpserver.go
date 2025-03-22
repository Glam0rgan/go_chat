package server

import (
	"go-im/common/discovery"
	"go-im/common/libnet"
	"go-im/connect/internal/svc"

	"github.com/alicebob/miniredis/v2/server"
	"github.com/zeromicro/go-zero/core/logx"
)

type TCPServer struct {
	svcCtx *svc.ServiceContext
	Server *server.Server
}

func (srv *TCPServer) HandleRequest() {

	for {
		session, err := srv.Server.Accept()
		if err != nil {
			panic(err)
		}

		cli := client.NewClient(srv.Server.Manager, session, srv.svcCtx.IMRpc)
		go srv.sessionLoop(cli)
	}
}

func (srv *TCPServer) sessionLoop(client *client.Clinet) {  

	message, err := client.Receive()
	if err != nil {
		logx.Errorf("[sessionLoop] client.Receive error: %v", err)
		_ = client.Close()
		return
	}

	// 登录校验
	err = client.Login(message)
	if err != nil {
		logx.Errorf("[sessionLoop] client.Login error: %v", err)
		_ = client.Close()
		return
	}

	go client.HeartBeat()

	for {
		message, err := client.Receive()
		if err != nil {
			logx.Errorf("[sessionLoop] client.Receive error: %v", err)
			_ = client.Close()
			return
		}
		err = client.HandlePackage(message)
		if err != nil {
			logx.Errorf("[sessionLoop] client.HandleMessage error: %v", err)
		}
	}
}

func NewSession(tcpServer *TCPServer, codec libnet.Codec, sendChanSize int) *Session {  
	sn := &Session{    
		tcpServer: tcpServer,    
		codec:     codec,    
		closeChan: make(chan int),  
	}  
	if sendChanSize > 0 {    
		sn.sendChan = make(chan libnet.Message, sendChanSize)    
		go sn.sendLoop()  
	}  
	return sn
}

func (srv *TCPServer) KqHeart() {  
	
	work := discovery.NewKqWorker(srv.svcCtx.Config.Etcd.Key, srv.svcCtx.Config.Etcd.Hosts, srv.svcCtx.Config.KqConf)  
	work.HeartBeat()
}