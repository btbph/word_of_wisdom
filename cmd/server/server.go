package main

import (
	"fmt"
	config "github.com/btbph/word_of_wisdom/internal/config/server"
	"github.com/btbph/word_of_wisdom/internal/model/server"
	"net"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server is running!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			// log error
			continue
		}
		client := server.NewClient(conn, cfg)
		go client.HandleRequests()
	}
}
