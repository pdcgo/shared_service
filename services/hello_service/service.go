package hello_service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
)

type helloServiceImpl struct{}

// Hello implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) Hello(context.Context, *connect.Request[access_iface.HelloRequest]) (*connect.Response[access_iface.HelloResponse], error) {
	panic("unimplemented")
}

func NewHelloService() *helloServiceImpl {
	return &helloServiceImpl{}
}
