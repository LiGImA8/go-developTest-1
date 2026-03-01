package rpc

import (
	"context"

	"google.golang.org/grpc"
)

type OrderServiceServer interface {
	PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error)
}

type OrderServiceClient struct{ cc *grpc.ClientConn }

func NewOrderServiceClient(cc *grpc.ClientConn) *OrderServiceClient {
	return &OrderServiceClient{cc: cc}
}

func (c *OrderServiceClient) PlaceOrder(ctx context.Context, in *PlaceOrderRequest) (*PlaceOrderResponse, error) {
	out := new(PlaceOrderResponse)
	err := c.cc.Invoke(ctx, "/"+OrderServiceName+"/PlaceOrder", in, out)
	return out, err
}

func RegisterOrderServiceServer(s *grpc.Server, srv OrderServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: OrderServiceName,
		HandlerType: (*OrderServiceServer)(nil),
		Methods: []grpc.MethodDesc{{
			MethodName: "PlaceOrder",
			Handler: func(si any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
				in := new(PlaceOrderRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.PlaceOrder(ctx, in)
				}
				info := &grpc.UnaryServerInfo{Server: si, FullMethod: "/" + OrderServiceName + "/PlaceOrder"}
				handler := func(ctx context.Context, req any) (any, error) { return srv.PlaceOrder(ctx, req.(*PlaceOrderRequest)) }
				return interceptor(ctx, in, info, handler)
			},
		}},
	}, srv)
}
