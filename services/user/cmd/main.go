package main

import (
	"net"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/pkg/logging"
	"minigate/pkg/rpc"
	"minigate/services/user/internal"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func main() {
	logger := logging.New("user-service")
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	encoding.RegisterCodec(rpc.JSONCodec{})
	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		logger.Fatal("open mysql failed", zap.Error(err))
	}
	defer mysql.Close()

	lis, err := net.Listen("tcp", config.GetEnv("USER_GRPC_ADDR", ":50051"))
	if err != nil {
		logger.Fatal("listen failed", zap.Error(err))
	}
	server := grpc.NewServer()
	rpc.RegisterUserServiceServer(server, internal.NewService(mysql))
	logger.Info("user service running", zap.String("addr", lis.Addr().String()))
	if err := server.Serve(lis); err != nil {
		logger.Fatal("grpc server stopped", zap.Error(err))
	}
}
