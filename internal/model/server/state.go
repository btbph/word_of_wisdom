package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/clock"
	"github.com/btbph/word_of_wisdom/internal/decode"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/btbph/word_of_wisdom/internal/dto/response"
	"github.com/btbph/word_of_wisdom/internal/hashcash"
	"github.com/btbph/word_of_wisdom/internal/usecase"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"strings"
	"time"
)

type Repo interface {
	SetChallengeInfo(ctx context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error
	RemoveChallengeInfo(ctx context.Context, ID uuid.UUID) error
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

func (s StandBy) Handle(ctx context.Context, connection ClientInterface, data io.Reader) ([]byte, error) {
	req, err := decode.JsonFromReader[request.RequestChallenge](data)
	if err != nil {
		s.logger.Error("failed to decode request challenge request", "error", err)
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	s.logger.Info("request for challenge recieved")

	if req.Type != dto.RequestChallenge {
		s.logger.Warn("expected request challenge request")
		return nil, errors.New("expect request challenge")
	}

	var (
		zeroBits   = connection.Config().Challenge.ZeroBits
		saltLength = connection.Config().Challenge.SaltLength
	)

	uc := usecase.NewInitChallenge(s.repo, s.logger)
	if err = uc.Init(ctx, connection.ClientID(), dto.NewChallengeInfo(zeroBits, saltLength)); err != nil {
		s.logger.Error("failed to set challenge info", "error", err)
		return nil, fmt.Errorf("failed to set challenge info: %w", err)
	}

	connection.SetState(NewWaitingForSolution(s.repo, s.logger))

	res := response.NewRequestChallengeResponse(zeroBits, saltLength)
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

func (s WaitingForSolution) Handle(ctx context.Context, connection ClientInterface, data io.Reader) ([]byte, error) {
	req, err := decode.JsonFromReader[request.SolutionProvided](data)
	if err != nil {
		s.logger.Error("failed to decode solution provided request", "error", err)
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	s.logger.Info("solution has been recieved")

	if req.Type != dto.SolutionProvided {
		s.logger.Warn("expected solution provided request request")
		return nil, errors.New("wrong request type")
	}

	challengeInfo, err := s.repo.ChallengeInfo(ctx, connection.ClientID())
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

	uc := usecase.NewCheckSolution(s.repo, validator, s.logger)
	quote, err := uc.Check(ctx, connection.ClientID(), req.Solution)
	if err != nil {
		return nil, err
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

func (s Finished) Handle(_ context.Context, _ ClientInterface, _ io.Reader) ([]byte, error) {
	return nil, nil
}
