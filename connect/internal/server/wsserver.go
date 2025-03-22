package server

import (
	"go-im/common/discovery"
	"go-im/common/socketio"
	"go-im/common/token"
	"go-im/connect/client"
	"go-im/connect/internal/svc"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WSServer struct {
	svcCtx *svc.ServiceContext
	Server *socketio.Server
}

func NewWSServer(svcCtx *svc.ServiceContext) *WSServer {
	return &WSServer{svcCtx: svcCtx}
}

func (ws *WSServer) Start() {
	err := http.ListenAndServe(ws.Server.Address, nil)
	if err != nil {
		panic(err)
	}
}

func (ws *WSServer) HandleRequest(w http.ResponseWriter, r *http.Request) {

	tk := r.FormValue("token")
	err := token.TokenValid(tk)

	if err != nil {
		logx.Infof("websocket token error")
		return
	}

	var upGrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
	}
	//cross origin domain support
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	session , err := ws.Server.Accept(conn)
	if err != nil {
		panic(err)
	}
	cli := client.NewClient(ws.Server.Manager, session, ws.svcCtx.IMRpc)
	ws.sessionLoop(cli)
}

// Fix
func (ws *WSServer) sessionLoop(client *client.Client) {
	/*
	for {
		message, err = client.Receive()
		if err != nil {
			logx.Errorf("[ws:sessionLoop] client.Receive error: %v", err)
			_ = client.Close()
			return
		}
		err = client.HandlePackage(message)
		if err != nil {
			logx.Errorf("[ws:sessionLoop] client.HandleMessage error: %v", err)
		}
	}
	*/
}

func (srv *WSServer) KqHeart() {  
	
	work := discovery.NewKqWorker(srv.svcCtx.Config.Etcd.Key, srv.svcCtx.Config.Etcd.Hosts, srv.svcCtx.Config.KqConf)  
	work.HeartBeat()
}