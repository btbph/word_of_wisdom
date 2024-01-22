package serializer

import (
	"errors"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/decode"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"io"
	"log/slog"
)

type RequestChallenge struct {
	logger *slog.Logger
}

func NewRequestChallenge(logger *slog.Logger) *RequestChallenge {
	return &RequestChallenge{logger: logger}
}

func (s *RequestChallenge) Serialize(data io.Reader) (request.RequestChallenge, error) {
	req, err := decode.JsonFromReader[request.RequestChallenge](data)
	if err != nil {
		s.logger.Error("failed to decode request challenge request", "error", err)
		return request.RequestChallenge{}, fmt.Errorf("failed to decode request: %w", err)
	}

	if req.Type != dto.RequestChallenge {
		s.logger.Warn("expected request challenge request")
		return request.RequestChallenge{}, errors.New("expect request challenge")
	}

	return req, nil
}
