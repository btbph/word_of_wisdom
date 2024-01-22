package repo

import (
	"context"
	"errors"
	"github.com/btbph/word_of_wisdom/internal/data_structures"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/google/uuid"
	"math/rand"
)

var ChallengeInfoNotFound = errors.New("challenge info not found")

type Repo struct {
	challengesInfo      *data_structures.ConcurrentMap[uuid.UUID, dto.ChallengeInfo]
	usedHashcashStrings *data_structures.ConcurrentSet[string]
	quotes              []string
}

func NewRepo() *Repo {
	return &Repo{
		challengesInfo:      data_structures.NewConcurrentMap[uuid.UUID, dto.ChallengeInfo](),
		usedHashcashStrings: data_structures.NewConcurrentSet[string](),
		quotes: []string{
			"Guard well your thoughts when alone and your words when accompanied.",
			"I like to listen. I have learned a great deal from listening carefully. Most people never listen.",
			"I think, that if the world were a bit more like ComicCon, it would be a better place.",
			"We must believe that we are gifted for something, and that this thing, at whatever cost, must be attained.",
			"The older I get, the greater power I seem to have to help the world; I am like a snowball - the further I am rolled the more I gain.",
			"Knowledge is love and light and vision",
		},
	}
}

func (r *Repo) SetChallengeInfo(_ context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error {
	r.challengesInfo.Insert(ID, challengeInfo)
	return nil
}

func (r *Repo) RemoveChallengeInfo(_ context.Context, ID uuid.UUID) error {
	r.challengesInfo.Delete(ID)
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
	index := rand.Intn(len(r.quotes))
	return r.quotes[index], nil
}
