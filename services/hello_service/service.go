package hello_service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/access_iface/v1"
	"github.com/pdcgo/shared/custom_connect"
)

type helloServiceImpl struct{}

// HelloBidiStream implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) HelloBidiStream(context.Context, *connect.BidiStream[access_iface.HelloBidiStreamRequest, access_iface.HelloBidiStreamResponse]) error {
	panic("unimplemented")
}

// HelloClientStream implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) HelloClientStream(context.Context, *connect.ClientStream[access_iface.HelloClientStreamRequest]) (*connect.Response[access_iface.HelloClientStreamResponse], error) {
	panic("unimplemented")
}

// HelloServerStream implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) HelloServerStream(context.Context, *connect.Request[access_iface.HelloServerStreamRequest], *connect.ServerStream[access_iface.HelloServerStreamResponse]) error {
	panic("unimplemented")
}

// Hello implements access_ifaceconnect.HelloServiceHandler.
func (h *helloServiceImpl) Hello(
	ctx context.Context,
	req *connect.Request[access_iface.HelloRequest]) (*connect.Response[access_iface.HelloResponse], error) {
	source, err := custom_connect.GetRequestSource(ctx)
	if err != nil {
		return nil, err
	}

	return &connect.Response[access_iface.HelloResponse]{
		Msg: &access_iface.HelloResponse{
			Source: source,
		},
	}, nil
}

func NewHelloService() *helloServiceImpl {
	return &helloServiceImpl{}
}
