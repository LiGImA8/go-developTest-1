package rpc

import (
	"context"

	"google.golang.org/grpc"
)

const (
	UserServiceName  = "minigate.v1.UserService"
	OrderServiceName = "minigate.v1.OrderService"
)

type UserServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
}

type UserServiceClient struct{ cc *grpc.ClientConn }

func NewUserServiceClient(cc *grpc.ClientConn) *UserServiceClient { return &UserServiceClient{cc: cc} }

func (c *UserServiceClient) Login(ctx context.Context, in *LoginRequest) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/"+UserServiceName+"/Login", in, out)
	return out, err
}

func (c *UserServiceClient) ValidateToken(ctx context.Context, in *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	out := new(ValidateTokenResponse)
	err := c.cc.Invoke(ctx, "/"+UserServiceName+"/ValidateToken", in, out)
	return out, err
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: UserServiceName,
		HandlerType: (*UserServiceServer)(nil),
		Methods: []grpc.MethodDesc{{
			MethodName: "Login",
			Handler: func(si any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
				in := new(LoginRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.Login(ctx, in)
				}
				info := &grpc.UnaryServerInfo{Server: si, FullMethod: "/" + UserServiceName + "/Login"}
				handler := func(ctx context.Context, req any) (any, error) { return srv.Login(ctx, req.(*LoginRequest)) }
				return interceptor(ctx, in, info, handler)
			},
		}, {
			MethodName: "ValidateToken",
			Handler: func(si any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
				in := new(ValidateTokenRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.ValidateToken(ctx, in)
				}
				info := &grpc.UnaryServerInfo{Server: si, FullMethod: "/" + UserServiceName + "/ValidateToken"}
				handler := func(ctx context.Context, req any) (any, error) {
					return srv.ValidateToken(ctx, req.(*ValidateTokenRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		}},
	}, srv)
}
