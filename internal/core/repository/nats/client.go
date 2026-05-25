package core_nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	Conn *nats.Conn
	JS   jetstream.JetStream
}

func New(ctx context.Context, cfg NatsConfig) (*Client, error) {
	nc, err := nats.Connect(
		cfg.URL,
		nats.Name(cfg.Name),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(cfg.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create jetstream client: %w", err)
	}

	return &Client{
		Conn: nc,
		JS:   js,
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.Conn == nil {
		return nil
	}

	return c.Conn.Drain()
}
