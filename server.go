package socket

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"go-im/common/libnet"
)

type Server struct {
	Name         string
	Manager      *libnet.Manager
	Listener     net.Listener
	SendChanSize int
}

func NewServer(name string, l net.Listener, sendChanSize int) *Server {
	return &Server{
		Name:         name,
		Manager:      libnet.NewManager(name),
		Listener:     l,
		SendChanSize: sendChanSize,
	}
}

func (s *Server) Accept() (*libnet.Session, error) {
	var tempDelay time.Duration
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if maxDelay := 1 * time.Second; tempDelay > maxDelay {
					tempDelay = maxDelay
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}

		var conn *websocket.connect

		return libnet.NewSession(s.Manager, s.SendChanSize), nil
	}
}

func (s *Server) Close() {
	s.Listener.Close()
	s.Manager.Close()
}

func NewServe(name, address string, sendChanSize int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewServer(name, listener, sendChanSize), nil
}