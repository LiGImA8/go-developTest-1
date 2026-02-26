package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"minigate/pkg/logger"

	"github.com/pebbe/zmq4"
)

type Consumer struct {
	db     *sql.DB
	socket *zmq4.Socket
}

func NewConsumer(db *sql.DB, endpoint string) (*Consumer, error) {
	socket, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		return nil, err
	}
	if err := socket.Bind(endpoint); err != nil {
		return nil, err
	}
	return &Consumer{db: db, socket: socket}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.socket.RecvBytes(0)
			if err != nil {
				log.Println("recv log message error:", err)
				continue
			}
			var evt logger.Event
			if err := json.Unmarshal(msg, &evt); err != nil {
				log.Println("invalid event payload", err)
				continue
			}
			_, err = c.db.ExecContext(ctx, `INSERT INTO logs(service, action, message, created_at) VALUES(?,?,?,?)`, evt.Service, evt.Action, evt.Message, evt.CreatedAt)
			if err != nil {
				log.Println("insert log failed", err)
			}
		}
	}
}

func (c *Consumer) Close() error { return c.socket.Close() }
