package logic

import (
	"context"

	"go-im/dao"
	"go-im/imrpc/imrpc"
	"go-im/imrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostMessageLogic) PostMessage(in *imrpc.PostMsg) (*imrpc.PostResponse, error) {
	
	u := dao.User{}
	userData := u.CheckHaveUserName(in.ToUserName)


	return &imrpc.PostResponse{}, nil
}
