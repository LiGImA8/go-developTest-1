package internal

import (
	"context"
	"database/sql"
	"encoding/json"

	zmq4 "github.com/go-zeromq/zmq4"
	"go.uber.org/zap"
	"minigate/pkg/logger"
)

type Consumer struct {
	db     *sql.DB
	socket zmq4.Socket
}

func NewConsumer(db *sql.DB, endpoint string) (*Consumer, error) {
	socket := zmq4.NewPull(context.Background())
	if err := socket.Listen(endpoint); err != nil {
		return nil, err
	}
	return &Consumer{db: db, socket: socket}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.socket.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				zap.L().Error("recv log message error", zap.Error(err))
				continue
			}
		}

		var evt logger.Event
		if err := json.Unmarshal(msg.Bytes(), &evt); err != nil {
			zap.L().Warn("invalid event payload", zap.Error(err))
			continue
		}
		_, err = c.db.ExecContext(ctx, `INSERT INTO logs(service, action, message, created_at) VALUES(?,?,?,?)`, evt.Service, evt.Action, evt.Message, evt.CreatedAt)
		if err != nil {
			zap.L().Error("insert log failed", zap.Error(err), zap.String("service", evt.Service), zap.String("action", evt.Action))
			continue
		}
		zap.L().Info("log persisted", zap.String("service", evt.Service), zap.String("action", evt.Action))
	}
}

func (c *Consumer) Close() error { return c.socket.Close() }
