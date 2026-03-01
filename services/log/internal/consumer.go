package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	zmq4 "github.com/go-zeromq/zmq4"
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
				log.Println("recv log message error:", err)
				continue
			}
		}

		var evt logger.Event
		if err := json.Unmarshal(msg.Bytes(), &evt); err != nil {
			log.Println("invalid event payload", err)
			continue
		}
		_, err = c.db.ExecContext(ctx, `INSERT INTO logs(service, action, message, created_at) VALUES(?,?,?,?)`, evt.Service, evt.Action, evt.Message, evt.CreatedAt)
		if err != nil {
			log.Println("insert log failed", err)
		}
	}
}

func (c *Consumer) Close() error { return c.socket.Close() }
