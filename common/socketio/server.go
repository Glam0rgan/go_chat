package socketio

import (
	"go-im/common/libnet"

	"github.com/gorilla/websocket"
)

type Server struct {
	Name         string
	Address      string
	Manager      *libnet.Manager
	SendChanSize int
}

func NewServe(name, address string, sendChanSize int) (*Server, error) {
	return &Server{
		Name:         name,
		Address:      address,
		Manager:      libnet.NewManager(name),
		SendChanSize: sendChanSize,
	}, nil
}

func (s *Server) Accept(conn *websocket.Conn) (*libnet.Session, error) {
	return libnet.NewSession(
		s.Manager,
		conn,
		s.SendChanSize,
	), nil
}