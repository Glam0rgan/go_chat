package logic

import (
	"context"
	"encoding/json"
	"go-im/common/libnet"
	"go-im/common/socketio"
	"go-im/connect/internal/svc"
	"go-im/imrpc/imrpc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type MqLogic struct {
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	wsServer *socketio.Server
	logx.Logger
}

type ConsumeMsg struct {
	Msg       string    `" json:"msg,omitempty"`
	Token     string    `" json:"Token,omitempty"`
	UserId    string    `" json:"ToUserName,omitempty"`
}

func NewMqLogic(ctx context.Context, svcCtx *svc.ServiceContext, wsSrv *socketio.Server) *MqLogic {
	return &MqLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		wsServer: wsSrv,
		Logger:   logx.WithContext(ctx),
	}
}

func (l *MqLogic) Consume(ctx context.Context , _, val string) error {

	logx.Infof("receive message")

	var msg imrpc.PostMsg
	err := json.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("[Consume] proto.Unmarshal val: %s error: %v", val, err)
		return err
	}
	logx.Infof("[Consume] succ msg: %+v body: %s", msg, msg.Msg)

	sess := l.wsServer.Manager.GetTokenSessions(msg.ToToken)
	if sess == nil {
		logx.Errorf("[Consume] session not found, msg: %+v", &msg)
		return nil
	}

	for _, ses := range sess {
		err = ses.Send(makeMessage(&msg))
		if err != nil {
			logx.Errorf("[Consume] session send error, msg: %+v, err: %v", &msg, err)
		}
	}

	return err
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext, wsSrv *socketio.Server) []service.Service {
	logx.Infof("create new consume")
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConf, NewMqLogic(ctx, svcCtx, wsSrv)),
	}
}

func makeMessage(msg *imrpc.PostMsg) libnet.Message {
	var message libnet.Message
	
	message.Body = []byte(msg.Msg)
	return message
}