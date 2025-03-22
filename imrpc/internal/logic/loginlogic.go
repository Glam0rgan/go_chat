package logic

import (
	"context"
	"errors"
	"go-im/common/token"
	"go-im/dao"

	"go-im/imrpc/imrpc"
	"go-im/imrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func LoginCheck(username, password string) (string, error) {
	var err error

	u := dao.User{}
	uData := u.CheckHaveUserName(username)

	if (uData.Id == 0) || (password != uData.Password) {
		err = errors.New("no this user or password error!")
		return "", err
	}

	token, err := token.GenerateToken(uData.Id)
	if err != nil {
		return "", err
	}
	return token, err
}

func (l *LoginLogic) Login(in *imrpc.LoginRequest) (*imrpc.LoginResponse, error) {
	
	// todo: jwt token
	token, err := LoginCheck(in.Username, in.Password)
	if err != nil {
		return nil, err
	}

	return &imrpc.LoginResponse{Token: token}, nil
}
