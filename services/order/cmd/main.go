package main

import (
	"net"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/pkg/logger"
	"minigate/pkg/logging"
	"minigate/pkg/rpc"
	"minigate/services/order/internal"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

func main() {
	log := logging.New("order-service")
	defer log.Sync()
	zap.ReplaceGlobals(log)

	encoding.RegisterCodec(rpc.JSONCodec{})
	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		log.Fatal("open mysql failed", zap.Error(err))
	}
	defer mysql.Close()

	userAddr := config.GetEnv("USER_SERVICE_ADDR", "localhost:50051")
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(rpc.JSONCodec{})))
	if err != nil {
		log.Fatal("connect user-service failed", zap.Error(err), zap.String("addr", userAddr))
	}
	defer userConn.Close()

	zmqEndpoint := config.GetEnv("ZMQ_LOG_ENDPOINT", "tcp://localhost:5557")
	publisher, err := logger.NewPublisher(zmqEndpoint)
	if err != nil {
		log.Fatal("connect log-service zmq failed", zap.Error(err), zap.String("endpoint", zmqEndpoint))
	}
	defer publisher.Close()

	lis, err := net.Listen("tcp", config.GetEnv("ORDER_GRPC_ADDR", ":50052"))
	if err != nil {
		log.Fatal("listen failed", zap.Error(err))
	}
	server := grpc.NewServer()
	rpc.RegisterOrderServiceServer(server, internal.NewService(mysql, rpc.NewUserServiceClient(userConn), publisher))
	log.Info("order service running", zap.String("addr", lis.Addr().String()), zap.String("user_service_addr", userAddr), zap.String("zmq_endpoint", zmqEndpoint))
	if err := server.Serve(lis); err != nil {
		log.Fatal("grpc server stopped", zap.Error(err))
	}
}
