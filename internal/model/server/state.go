package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/clock"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/btbph/word_of_wisdom/internal/dto/response"
	"github.com/btbph/word_of_wisdom/internal/hashcash"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"strings"
	"time"
)

type Repo interface {
	SetChallengeInfo(ctx context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error
	ChallengeInfo(ctx context.Context, ID uuid.UUID) (dto.ChallengeInfo, error)
	CheckSolutionPresence(ctx context.Context, solution string) (bool, error)
	SaveSolution(ctx context.Context, solution string) error
	RandomQuote(ctx context.Context) (string, error)
}

type StandBy struct {
	repo   Repo
	logger *slog.Logger
}

func NewStandBy(repo Repo, logger *slog.Logger) *StandBy {
	return &StandBy{
		repo:   repo,
		logger: logger,
	}
}

func (s StandBy) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	req := request.RequestChallenge{}
	if err := json.NewDecoder(data).Decode(&req); err != nil {
		s.logger.Error("failed to decode request challenge request", "error", err)
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	s.logger.Info("request for challenge recieved")

	if req.Type != dto.RequestChallenge {
		s.logger.Warn("expected request challenge request")
		return nil, errors.New("expect request challenge")
	}

	connection.SetState(NewWaitingForSolution(s.repo, s.logger))

	var (
		zeroBits   = connection.Config().Challenge.ZeroBits
		saltLength = connection.Config().Challenge.SaltLength
	)
	res := response.NewRequestChallengeResponse(zeroBits, saltLength)
	if err := s.repo.SetChallengeInfo(context.TODO(), connection.ClientID(), dto.NewChallengeInfo(zeroBits, saltLength)); err != nil {
		s.logger.Error("failed to set challenge info", "error", err)
		return nil, fmt.Errorf("failed to set challenge info: %w", err)
	}

	return json.Marshal(res)
}

type WaitingForSolution struct {
	repo   Repo
	logger *slog.Logger
}

func NewWaitingForSolution(repo Repo, logger *slog.Logger) *WaitingForSolution {
	return &WaitingForSolution{
		repo:   repo,
		logger: logger,
	}
}

func (s WaitingForSolution) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	req := request.SolutionProvided{}
	if err := json.NewDecoder(data).Decode(&req); err != nil {
		s.logger.Error("failed to decode solution provided request", "error", err)
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	s.logger.Info("solution has been recieved")

	if req.Type != dto.SolutionProvided {
		s.logger.Warn("expected solution provided request request")
		return nil, errors.New("wrong request type")
	}

	challengeInfo, err := s.repo.ChallengeInfo(context.TODO(), connection.ClientID())
	if err != nil {
		s.logger.Error("failed to get challenge info", "error", err)
		return nil, fmt.Errorf("failed to get challenge info: %w", err)
	}

	expireDate := time.Duration(connection.Config().Challenge.ExpireDateInHours) * time.Hour
	validator := hashcash.NewValidator(
		challengeInfo,
		clock.New(),
		expireDate,
		sha256.New(),
		s.expectedResource(connection.Config().Server.Address),
	)
	if err = validator.Validate(req.Solution); err != nil {
		s.logger.Warn("solution validation failed", "error", err)
		return nil, fmt.Errorf("solution validation failed: %w", err)
	}

	present, err := s.repo.CheckSolutionPresence(context.TODO(), req.Solution)
	if err != nil {
		s.logger.Error("failed to check solution presence", "error", err)
		return nil, fmt.Errorf("failed to check solution presence: %w", err)
	}

	if present {
		s.logger.Warn("current solution already presents", "solution", req.Solution)
		return nil, errors.New("current solution already presents")
	}

	if err = s.repo.SaveSolution(context.TODO(), req.Solution); err != nil {
		s.logger.Error("failed to save solution", "error", err)
		return nil, errors.New("failed to save solution")
	}

	quote, err := s.repo.RandomQuote(context.TODO())
	if err != nil {
		s.logger.Error("failed to get random quote", "error", err)
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	connection.SetState(Finished{})
	connection.Close()

	res := response.NewSolutionProvided(quote)
	return json.Marshal(res)
}

func (s WaitingForSolution) expectedResource(address string) string {
	return strings.Split(address, ":")[0]
}

type Finished struct{}

func (s Finished) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	return nil, nil
}
