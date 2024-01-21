package repo

import (
	"context"
	"errors"
	"github.com/btbph/word_of_wisdom/internal/data_structures"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/google/uuid"
)

var ChallengeInfoNotFound = errors.New("challenge info not found")

type Repo struct {
	challengesInfo      *data_structures.ConcurrentMap[uuid.UUID, dto.ChallengeInfo]
	usedHashcashStrings *data_structures.ConcurrentSet[string]
}

func NewRepo() *Repo {
	return &Repo{
		challengesInfo:      data_structures.NewConcurrentMap[uuid.UUID, dto.ChallengeInfo](),
		usedHashcashStrings: data_structures.NewConcurrentSet[string](),
	}
}

func (r *Repo) SetChallengeInfo(_ context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error {
	r.challengesInfo.Insert(ID, challengeInfo)
	return nil
}

func (r *Repo) ChallengeInfo(_ context.Context, ID uuid.UUID) (dto.ChallengeInfo, error) {
	res, ok := r.challengesInfo.Get(ID)
	if !ok {
		return dto.ChallengeInfo{}, ChallengeInfoNotFound
	}

	return res, nil
}

func (r *Repo) CheckSolutionPresence(_ context.Context, solution string) (bool, error) {
	ok := r.usedHashcashStrings.Exist(solution)
	return ok, nil
}

func (r *Repo) SaveSolution(_ context.Context, solution string) error {
	r.usedHashcashStrings.Insert(solution)
	return nil
}

func (r *Repo) RandomQuote(_ context.Context) (string, error) {
	return "Random quote!", nil
}
