package main

import (
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/model/server"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
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
		client := server.NewClient(conn)
		go client.HandleRequests()
	}
}
