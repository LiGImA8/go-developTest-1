package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pebbe/zmq4"
)

type Event struct {
	Service   string    `json:"service"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Publisher struct {
	socket *zmq4.Socket
}

func NewPublisher(endpoint string) (*Publisher, error) {
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		return nil, err
	}
	if err := socket.Connect(endpoint); err != nil {
		return nil, err
	}
	return &Publisher{socket: socket}, nil
}

func (p *Publisher) Publish(ctx context.Context, event Event) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err = p.socket.SendBytes(payload, 0)
		if err != nil {
			return fmt.Errorf("send zmq: %w", err)
		}
		return nil
	}
}

func (p *Publisher) Close() error {
	return p.socket.Close()
}
