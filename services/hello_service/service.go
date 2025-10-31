package hello_service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/shared/custom_connect"
)

type helloServiceImpl struct{}

// Hello implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) Hello(
	ctx context.Context,
	req *connect.Request[access_iface.HelloRequest]) (*connect.Response[access_iface.HelloResponse], error) {
	return &connect.Response[access_iface.HelloResponse]{
		Msg: &access_iface.HelloResponse{
			Source: custom_connect.GetRequestSource(ctx),
		},
	}, nil
}

func NewHelloService() *helloServiceImpl {
	return &helloServiceImpl{}
}
