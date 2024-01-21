package main

import (
	config "github.com/btbph/word_of_wisdom/internal/config/server"
	"github.com/btbph/word_of_wisdom/internal/logger"
	"github.com/btbph/word_of_wisdom/internal/model/server"
	"net"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	l, err := logger.Init(cfg.Server.LogLevel)
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		l.Error("failed to establish tcp listener", "error", err)
		panic(err)
	}
	defer func() {
		if err = listener.Close(); err != nil {
			l.Error("failed to close listener", "error", err)
		}
	}()

	l.Info("server is running")

	for {
		conn, err := listener.Accept()
		if err != nil {
			l.Error("failed to accept connection", "error", err)
		}
		l.Info("connection is established")

		client := server.NewClient(conn, cfg, l)
		go client.HandleRequests()
	}
}
