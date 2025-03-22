package client

import (
	"context"
	"go-im/common/libnet"
	"go-im/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Session *libnet.Session
	Manager *libnet.Manager
	IMRpc   imrpcclient.Imrpc
	heartbeat chan *libnet.Message
}

func NewClient(manager *libnet.Manager, session *libnet.Session, imrpc imrpcclient.Imrpc) *Client {
	return &Client{
		Session:   session,
		Manager:   manager,
		IMRpc:     imrpc,
		heartbeat: make(chan *libnet.Message),
	}
}

func (c *Client) HandlePackage(msg *libnet.Message) error {
	
	req := makePostMessage(c.Session.Session().String(), msg)
	if req == nil {
		return nil
	}
	_, err := c.IMRpc.PostMessage(context.Background(), req)
	if err != nil {
		logx.Errorf("[HandlePackage] client.PostMessage error: %v", err)
	}

	return err
}

func makePostMessage(sessionId string, msg *libnet.Message) *imrpcclient.PostMsg {
	var postMessageReq imrpcclient.PostMsg
	err := proto.Unmarshal(msg.Body, &postMessageReq)
	if err != nil {
		logx.Errorf("[makePostMessage] proto.Unmarshal msg: %v error: %v", msg, err)
		return nil
	}
	
	postMessageReq.SessionId = sessionId

	return &postMessageReq
}

/*
func (c *Client) Receive() (*libnet.Message, error) {
	return c.Session.Receive()
}
*/

func (c *Client) Send(msg libnet.Message) error {
	return c.Session.Send(msg)
}