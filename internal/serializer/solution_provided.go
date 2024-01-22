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

type SolutionProvided struct {
	logger *slog.Logger
}

func NewSolutionProvided(logger *slog.Logger) *SolutionProvided {
	return &SolutionProvided{logger: logger}
}

func (s *SolutionProvided) Serialize(data io.Reader) (request.SolutionProvided, error) {
	req, err := decode.JsonFromReader[request.SolutionProvided](data)
	if err != nil {
		s.logger.Error("failed to decode solution provided request", "error", err)
		return request.SolutionProvided{}, fmt.Errorf("failed to decode request: %w", err)
	}

	if req.Type != dto.SolutionProvided {
		s.logger.Warn("expected solution provided request request")
		return request.SolutionProvided{}, errors.New("wrong request type")
	}

	return req, nil
}
