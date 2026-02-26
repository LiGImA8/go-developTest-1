package main

import (
	"log"
	"net"

	"minigate/pkg/config"
	"minigate/pkg/db"
	"minigate/pkg/rpc"
	"minigate/services/user/internal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func main() {
	encoding.RegisterCodec(rpc.JSONCodec{})
	mysql, err := db.Open(config.MustDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer mysql.Close()

	lis, err := net.Listen("tcp", config.GetEnv("USER_GRPC_ADDR", ":50051"))
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	rpc.RegisterUserServiceServer(server, internal.NewService(mysql))
	log.Println("user service running on", lis.Addr())
	log.Fatal(server.Serve(lis))
}
