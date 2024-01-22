package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/clock"
	config "github.com/btbph/word_of_wisdom/internal/config/client"
	"github.com/btbph/word_of_wisdom/internal/decode"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/btbph/word_of_wisdom/internal/dto/response"
	"github.com/btbph/word_of_wisdom/internal/hashcash"
	"github.com/btbph/word_of_wisdom/internal/logger"
	"log/slog"
	"net"
	"strings"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	l, err := logger.Init(cfg.Client.LogLevel)
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", cfg.Client.DestinationAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			l.Error("failed to close connection", "error", err)
		}
	}()

	if err = requestChallenge(conn, l); err != nil {
		return
	}

	challenge, err := recieveChallenge(conn, l)
	if err != nil {
		return
	}

	if err = sendChallengeSolution(conn, challenge, resource(cfg.Client.DestinationAddress), l); err != nil {
		return
	}

	quote, err := recieveQuote(conn, l)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Quote recieved: %s", quote))
}

func requestChallenge(conn net.Conn, logger *slog.Logger) error {
	bytes, err := request.MarshalRequestChallenge(request.NewRequestChallenge())
	if err != nil {
		logger.Error("failed to marshall request challenge request", "error", err)
		return fmt.Errorf("failed to marshall request challenge request: %w", err)
	}

	_, err = conn.Write(bytes)
	if err != nil {
		logger.Error("failed to write connection", "error", err)
		return fmt.Errorf("failed to write connection: %w", err)
	}

	return nil
}

func recieveChallenge(conn net.Conn, logger *slog.Logger) (response.RequestChallenge, error) {
	res, err := decode.JsonFromReader[response.RequestChallenge](conn)
	if err != nil {
		logger.Error("failed to decode reqest challenge response", "error", err)
		return response.RequestChallenge{}, fmt.Errorf("failed to decode reqest challenge response: %w", err)
	}
	return res, nil
}

func resource(destinationAddress string) string {
	return strings.Split(destinationAddress, ":")[0]
}

func sendChallengeSolution(conn net.Conn, challenge response.RequestChallenge, resource string, logger *slog.Logger) error {
	h := hashcash.NewHashcash(challenge.ZeroBits, challenge.SaltLength, clock.New(), sha256.New())
	req := request.NewSolutionProvided(h.Generate(resource))

	bytes, err := request.MarshalSolutionProvided(req)
	if err != nil {
		logger.Error("failed to marshall solution provided response", "error", err)
		return err
	}

	_, err = conn.Write(bytes)
	if err != nil {
		logger.Error("failed to write solution provided response", "error", err)
		return fmt.Errorf("failed to write solution provided response: %w", err)
	}

	return nil
}

func recieveQuote(conn net.Conn, logger *slog.Logger) (string, error) {
	res, err := decode.JsonFromReader[response.SolutionProvided](conn)
	if err != nil {
		logger.Error("failed to decode solution provided response", "error", err)
		return "", fmt.Errorf("failed to decode solution provided response: %w", err)
	}

	return res.Quote, nil
}
