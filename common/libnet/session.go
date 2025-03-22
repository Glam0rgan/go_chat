package libnet

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go-im/common/session"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	globalSessionId = uint64(rand.New(rand.NewSource(time.Now().Unix())).Uint32())
}

var (
	SessionClosedError  = errors.New("session Closed")
	SessionBlockedError = errors.New("session Blocked")

	globalSessionId uint64
)

type Session struct {
	id         uint64
	token      string
	conn       *websocket.Conn
	manager    *Manager
	sendChan   chan Message
	closeFlag  int32
	closeChan  chan int
	closeMutex sync.Mutex
}

func NewSession(manager *Manager,conn *websocket.Conn, sendChanSize int) *Session {
	s := &Session{
		manager:   manager,
		conn: 		conn,
		closeChan: make(chan int),
		id:        atomic.AddUint64(&globalSessionId, 1),
	}
	if sendChanSize > 0 {
		s.sendChan = make(chan Message, sendChanSize)
		go s.sendLoop()
	}

	return s
}

func (s *Session) Name() string {
	return s.manager.Name
}

func (s *Session) ID() uint64 {
	return s.id
}

func (s *Session) Token() string {
	return s.token
}

func (s *Session) Session() session.Session {
	return session.NewSession(s.manager.Name, s.token, s.id)
}

func (s *Session) SetToken(token string) {
	s.token = token
}

func (s *Session) sendLoop() {
	
	defer s.Close()

	for {
		message, ok := <- s.sendChan
		s.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
		if !ok {
			logx.Infof("SetWriteDeadline not ok")
			s.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := s.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			logx.Infof(" ch.conn.NextWriter err :%s  ", err.Error())
			return
		}
		logx.Infof("message write body:%s", message.Body)
		w.Write(message.Body)
		if err := w.Close(); err != nil {
			return
		}
	}
}


func(s *Session) Send(msg Message) error {
	if s.IsClosed() {
		return SessionClosedError
	}
	if s.sendChan == nil {
		s.sendChan = make(chan Message)
	}
	select {
	case s.sendChan <- msg:
		return nil
	default:
		return SessionBlockedError
	}
}


/*
func (s *Session) Receive() (*Message, error) {
	
}*/

func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

func (s *Session) Close() error {
	/*
	if atomic.CompareAndSwapInt32(&s.closeFlag, 0, 1) {
		err := s.codec.Close()
		close(s.closeChan)
		if s.manager != nil {
			s.manager.removeSession(s)
		}
		return err
	}
	
	*/
	return SessionClosedError
}