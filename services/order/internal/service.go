package internal

import (
	"context"
	"database/sql"
	"fmt"

	"minigate/pkg/logger"
	"minigate/pkg/rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	db         *sql.DB
	userClient *rpc.UserServiceClient
	publisher  *logger.Publisher
}

func NewService(db *sql.DB, userClient *rpc.UserServiceClient, publisher *logger.Publisher) *Service {
	return &Service{db: db, userClient: userClient, publisher: publisher}
}

func (s *Service) PlaceOrder(ctx context.Context, req *rpc.PlaceOrderRequest) (*rpc.PlaceOrderResponse, error) {
	if req.Quantity <= 0 || req.ItemName == "" {
		return nil, status.Error(codes.InvalidArgument, "item_name and positive quantity are required")
	}
	validResp, err := s.userClient.ValidateToken(ctx, &rpc.ValidateTokenRequest{Token: req.Token})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	if !validResp.Valid {
		return nil, status.Error(codes.Unauthenticated, "token invalid")
	}
	res, err := s.db.ExecContext(ctx, `INSERT INTO orders(user_id, item_name, quantity, status) VALUES(?,?,?,?)`, validResp.UserID, req.ItemName, req.Quantity, "CREATED")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orderID, _ := res.LastInsertId()
	_ = s.publisher.Publish(ctx, logger.Event{
		Service: "order-service",
		Action:  "create_order",
		Message: fmt.Sprintf("order=%d user=%d item=%s qty=%d", orderID, validResp.UserID, req.ItemName, req.Quantity),
	})
	return &rpc.PlaceOrderResponse{OrderID: orderID, Status: "CREATED"}, nil
}
