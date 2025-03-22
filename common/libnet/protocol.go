package libnet

import "net"

type Message struct {
	SeqId uint32 
	Body  []byte 
}

type Codec interface {
	Close() error
	Send(Message) error
}

type imCodec struct {
	conn net.Conn
}