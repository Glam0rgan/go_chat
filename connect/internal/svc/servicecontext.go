package svc

import (
	"go-im/connect/internal/config"
	"go-im/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpcclient.Imrpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.IMRpc)

	return &ServiceContext{
		Config: c,
		IMRpc:  imrpcclient.NewImrpc(client),
	}
}
