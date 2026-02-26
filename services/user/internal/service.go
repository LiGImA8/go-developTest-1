package internal

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"minigate/pkg/rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct{ db *sql.DB }

func NewService(db *sql.DB) *Service { return &Service{db: db} }

func (s *Service) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username/password required")
	}
	var userID int64
	var dbPass string
	err := s.db.QueryRowContext(ctx, `SELECT id, password FROM users WHERE username=?`, req.Username).Scan(&userID, &dbPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if dbPass != req.Password {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}
	token, err := randomToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO tokens(user_id, token, expires_at) VALUES(?,?,?)`, userID, token, time.Now().Add(1*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &rpc.LoginResponse{Token: token, UserID: userID}, nil
}

func (s *Service) ValidateToken(ctx context.Context, req *rpc.ValidateTokenRequest) (*rpc.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &rpc.ValidateTokenResponse{Valid: false}, nil
	}
	var userID int64
	var username string
	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, u.username
		FROM tokens t JOIN users u ON t.user_id=u.id
		WHERE t.token=? AND t.expires_at > NOW()
	`, req.Token).Scan(&userID, &username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &rpc.ValidateTokenResponse{Valid: false}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &rpc.ValidateTokenResponse{Valid: true, UserID: userID, Username: username}, nil
}

func randomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
