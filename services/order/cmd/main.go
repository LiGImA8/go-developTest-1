package main

import (
	"log"
	"net"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/pkg/logger"
	"minigate/pkg/rpc"
	"minigate/services/order/internal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

func main() {
	encoding.RegisterCodec(rpc.JSONCodec{})
	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	userConn, err := grpc.NewClient(config.GetEnv("USER_SERVICE_ADDR", "user-service:50051"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(rpc.JSONCodec{})))
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	publisher, err := logger.NewPublisher(config.GetEnv("ZMQ_LOG_ENDPOINT", "tcp://log-service:5557"))
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	lis, err := net.Listen("tcp", config.GetEnv("ORDER_GRPC_ADDR", ":50052"))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	rpc.RegisterOrderServiceServer(server, internal.NewService(mysql, rpc.NewUserServiceClient(userConn), publisher))
	log.Println("order service running on", lis.Addr())
	log.Fatal(server.Serve(lis))
}
