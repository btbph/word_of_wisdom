package server

import (
	"github.com/btbph/word_of_wisdom/internal/config/server"
	"github.com/btbph/word_of_wisdom/internal/repo"
	"github.com/google/uuid"
	"io"
	"net"
)

type Client struct {
	conn     net.Conn
	state    ConnectionState
	finished bool
	id       uuid.UUID
	cfg      *server.Config
}

func NewClient(conn net.Conn, cfg *server.Config) *Client {
	return &Client{
		conn:  conn,
		state: NewStandBy(repo.NewRepo()),
		id:    uuid.New(),
		cfg:   cfg,
	}
}

func (c *Client) SetState(state ConnectionState) {
	c.state = state
}

func (c *Client) Close() {
	c.finished = true
}

func (c *Client) HandleRequests() {
	defer c.conn.Close()
	for !c.finished {
		resp, err := c.state.Handle(c, c.conn)
		if err != nil {
			// log error
			return
		}

		if len(resp) > 0 {
			_, err = c.conn.Write(resp)
			if err != nil {
				// log error
				return
			}
		}
	}
}

func (c *Client) ClientID() uuid.UUID {
	return c.id
}

func (c *Client) Config() *server.Config {
	return c.cfg
}

type ClientInterface interface {
	SetState(state ConnectionState)
	Close()
	ClientID() uuid.UUID
	Config() *server.Config
}

type ConnectionState interface {
	Handle(connection ClientInterface, data io.Reader) ([]byte, error)
}
