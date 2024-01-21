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
	repo Repo
}

func NewStandBy(repo Repo) *StandBy {
	return &StandBy{repo: repo}
}

func (s StandBy) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	req := request.RequestChallenge{}
	if err := json.NewDecoder(data).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	if req.Type != dto.RequestChallenge {
		return nil, errors.New("expect request challenge")
	}

	connection.SetState(NewWaitingForSolution(s.repo))
	const (
		zeroBits   = 20
		saltLength = 8
	)
	res := response.NewRequestChallengeResponse(zeroBits, saltLength) // TODO: take from config
	if err := s.repo.SetChallengeInfo(context.TODO(), connection.ClientID(), dto.NewChallengeInfo(zeroBits, saltLength)); err != nil {
		return nil, fmt.Errorf("failed to set challenge info: %w", err)
	}

	return json.Marshal(res)
}

type WaitingForSolution struct {
	repo Repo
}

func NewWaitingForSolution(repo Repo) *WaitingForSolution {
	return &WaitingForSolution{repo: repo}
}

func (s WaitingForSolution) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	req := request.SolutionProvided{}
	if err := json.NewDecoder(data).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	if req.Type != dto.SolutionProvided {
		return nil, errors.New("wrong request type")
	}

	challengeInfo, err := s.repo.ChallengeInfo(context.TODO(), connection.ClientID())
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge info: %w", err)
	}

	expireDate := 28 * 24 * time.Hour
	expectedResource := "localhost"
	validator := hashcash.NewValidator(challengeInfo, clock.New(), expireDate, sha256.New(), expectedResource)
	if err = validator.Validate(req.Solution); err != nil {
		return nil, fmt.Errorf("solution validation failed: %w", err)
	}

	present, err := s.repo.CheckSolutionPresence(context.TODO(), req.Solution)
	if err != nil {
		return nil, fmt.Errorf("failed to check solution presents: %w", err)
	}

	if present {
		return nil, errors.New("current solution already presents")
	}

	if err = s.repo.SaveSolution(context.TODO(), req.Solution); err != nil {
		return nil, errors.New("failed to save solution")
	}

	quote, err := s.repo.RandomQuote(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	connection.SetState(Finished{})
	connection.Close()

	res := response.NewSolutionProvided(quote)
	return json.Marshal(res)
}

type Finished struct{}

func (s Finished) Handle(connection ClientInterface, data io.Reader) ([]byte, error) {
	return nil, nil
}
