package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	zmq4 "github.com/go-zeromq/zmq4"
)

type Event struct {
	Service   string    `json:"service"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Publisher struct {
	socket zmq4.Socket
}

func NewPublisher(endpoint string) (*Publisher, error) {
	socket := zmq4.NewPush(context.Background())
	if err := socket.Dial(endpoint); err != nil {
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
		err = p.socket.Send(zmq4.NewMsg(payload))
		if err != nil {
			return fmt.Errorf("send zmq: %w", err)
		}
		return nil
	}
}

func (p *Publisher) Close() error {
	return p.socket.Close()
}
