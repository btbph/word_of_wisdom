package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/clock"
	config "github.com/btbph/word_of_wisdom/internal/config/server"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/dto/response"
	"github.com/btbph/word_of_wisdom/internal/hashcash"
	"github.com/btbph/word_of_wisdom/internal/serializer"
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
	cfg    *config.Config
	logger *slog.Logger
}

func NewStandBy(
	repo Repo,
	cfg *config.Config,
	logger *slog.Logger,
) *StandBy {
	return &StandBy{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s StandBy) Handle(ctx context.Context, connection ClientInterface, data io.Reader) ([]byte, error) {
	ser := serializer.NewRequestChallenge(s.logger)
	_, err := ser.Serialize(data)
	if err != nil {
		return nil, err
	}

	s.logger.Info("request for challenge recieved")

	var (
		zeroBits   = s.cfg.Challenge.ZeroBits
		saltLength = s.cfg.Challenge.SaltLength
	)

	uc := usecase.NewInitChallenge(s.repo, s.logger)
	if err = uc.Init(ctx, connection.ClientID(), dto.NewChallengeInfo(zeroBits, saltLength)); err != nil {
		s.logger.Error("failed to set challenge info", "error", err)
		return nil, fmt.Errorf("failed to set challenge info: %w", err)
	}

	connection.SetState(NewWaitingForSolution(s.repo, s.cfg, s.logger))

	res := response.NewRequestChallengeResponse(zeroBits, saltLength)
	return json.Marshal(res)
}

type WaitingForSolution struct {
	repo   Repo
	cfg    *config.Config
	logger *slog.Logger
}

func NewWaitingForSolution(repo Repo, cfg *config.Config, logger *slog.Logger) *WaitingForSolution {
	return &WaitingForSolution{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s WaitingForSolution) Handle(ctx context.Context, connection ClientInterface, data io.Reader) ([]byte, error) {
	ser := serializer.NewSolutionProvided(s.logger)
	req, err := ser.Serialize(data)
	if err != nil {
		return nil, err
	}

	s.logger.Info("solution has been recieved")

	challengeInfo, err := s.repo.ChallengeInfo(ctx, connection.ClientID())
	if err != nil {
		s.logger.Error("failed to get challenge info", "error", err)
		return nil, fmt.Errorf("failed to get challenge info: %w", err)
	}

	expireDate := time.Duration(s.cfg.Challenge.ExpireDateInHours) * time.Hour
	validator := hashcash.NewValidator(
		challengeInfo,
		clock.New(),
		expireDate,
		sha256.New(),
		s.expectedResource(s.cfg.Server.Address),
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
