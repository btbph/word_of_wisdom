package server

import (
	"context"
	"github.com/btbph/word_of_wisdom/internal/config/server"
	"github.com/btbph/word_of_wisdom/internal/repo"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net"
)

type Client struct {
	conn     net.Conn
	state    ConnectionState
	finished bool
	id       uuid.UUID

	logger *slog.Logger
}

func NewClient(conn net.Conn, cfg *server.Config, logger *slog.Logger) *Client {
	id := uuid.New()
	l := logger.With("connectionID", id)
	return &Client{
		conn:   conn,
		state:  NewStandBy(repo.NewRepo(), cfg, l),
		id:     id,
		logger: l,
	}
}

func (c *Client) SetState(state ConnectionState) {
	c.state = state
}

func (c *Client) Close() {
	c.finished = true
}

func (c *Client) HandleRequests() {
	defer func() {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("failed to close connection", "error", err)
		}
	}()
	for !c.finished {
		resp, err := c.state.Handle(context.TODO(), c, c.conn)
		if err != nil {
			return
		}

		if len(resp) > 0 {
			_, err = c.conn.Write(resp)
			if err != nil {
				c.logger.Error("failed to write response", "error", err)
				return
			}
		}
	}
}

func (c *Client) ClientID() uuid.UUID {
	return c.id
}

type ClientInterface interface {
	SetState(state ConnectionState)
	Close()
	ClientID() uuid.UUID
}

type ConnectionState interface {
	Handle(ctx context.Context, connection ClientInterface, data io.Reader) ([]byte, error)
}
