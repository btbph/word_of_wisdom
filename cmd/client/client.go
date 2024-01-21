package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/clock"
	config "github.com/btbph/word_of_wisdom/internal/config/client"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/btbph/word_of_wisdom/internal/dto/response"
	"github.com/btbph/word_of_wisdom/internal/hashcash"
	"net"
	"strings"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", cfg.Client.DestinationAddress)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err = requestChallenge(conn); err != nil {
		// log error
		return
	}

	challenge, err := recieveChallenge(conn)
	if err != nil {
		return
	}

	if err = sendChallengeSolution(conn, challenge, resource(cfg.Client.DestinationAddress)); err != nil {
		return
	}

	quote, err := recieveQuote(conn)
	if err != nil {
		return
	}

	fmt.Printf("Quote recieved: %s\n", quote)
}

func requestChallenge(conn net.Conn) error {
	req := request.NewRequestChallenge()
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func recieveChallenge(conn net.Conn) (response.RequestChallenge, error) {
	res := response.RequestChallenge{}
	if err := json.NewDecoder(conn).Decode(&res); err != nil {
		return response.RequestChallenge{}, err
	}

	return res, nil
}

func resource(destinationAddress string) string {
	return strings.Split(destinationAddress, ":")[0]
}

func sendChallengeSolution(conn net.Conn, challenge response.RequestChallenge, resource string) error {
	h := hashcash.NewHashcash(challenge.ZeroBits, challenge.SaltLength, clock.New(), sha256.New())
	req := request.NewSolutionProvided(h.Generate(resource))

	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)
	return err
}

func recieveQuote(conn net.Conn) (string, error) {
	res := response.SolutionProvided{}
	if err := json.NewDecoder(conn).Decode(&res); err != nil {
		return "", err
	}

	return res.Quote, nil
}
