package logic

import (
	"context"

	"go-im/connect/connect"
	"go-im/connect/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *connect.Request) (*connect.Response, error) {
	// todo: add your logic here and delete this line

	return &connect.Response{}, nil
}
