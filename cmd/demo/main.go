package main

import (
	"context"
	"log"
	"time"

	"minigate/pkg/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(rpc.JSONCodec{})))
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	orderConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.ForceCodec(rpc.JSONCodec{})))
	if err != nil {
		log.Fatal(err)
	}
	defer orderConn.Close()

	userClient := rpc.NewUserServiceClient(userConn)
	orderClient := rpc.NewOrderServiceClient(orderConn)

	loginResp, err := userClient.Login(ctx, &rpc.LoginRequest{Username: "demo", Password: "demo123"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("login ok user_id=%d token=%s", loginResp.UserID, loginResp.Token)

	orderResp, err := orderClient.PlaceOrder(ctx, &rpc.PlaceOrderRequest{Token: loginResp.Token, ItemName: "book", Quantity: 1})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("order ok id=%d status=%s", orderResp.OrderID, orderResp.Status)
}
