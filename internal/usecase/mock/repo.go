package mock

import (
	"context"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type Repo struct {
	mock.Mock
}

func (r *Repo) SetChallengeInfo(_ context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error {
	return r.Called(ID, challengeInfo).Error(0)
}

func (r *Repo) RemoveChallengeInfo(_ context.Context, ID uuid.UUID) error {
	return r.Called(ID).Error(0)
}

func (r *Repo) ChallengeInfo(_ context.Context, ID uuid.UUID) (dto.ChallengeInfo, error) {
	args := r.Called(ID)
	if err := args.Error(1); err != nil {
		return dto.ChallengeInfo{}, err
	}

	return args.Get(0).(dto.ChallengeInfo), nil
}

func (r *Repo) CheckSolutionPresence(_ context.Context, solution string) (bool, error) {
	args := r.Called(solution)
	if err := args.Error(1); err != nil {
		return false, err
	}

	return args.Bool(0), nil
}

func (r *Repo) SaveSolution(_ context.Context, solution string) error {
	return r.Called(solution).Error(0)
}

func (r *Repo) RandomQuote(_ context.Context) (string, error) {
	args := r.Called()
	if err := args.Error(1); err != nil {
		return "", err
	}

	return args.String(0), nil
}
