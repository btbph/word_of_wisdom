package usecase

import (
	"context"
	"errors"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/usecase/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"testing"
	"time"
)

type InitChallengeSuite struct {
	suite.Suite

	uc   *InitChallenge
	repo *mock.Repo

	logger *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func TestInitChallengeSuite(t *testing.T) {
	suite.Run(t, &InitChallengeSuite{})
}

func (s *InitChallengeSuite) SetupSuite() {
	s.logger = slog.Default()
}

func (s *InitChallengeSuite) SetupTest() {
	s.repo = &mock.Repo{}
	s.uc = NewInitChallenge(s.repo, s.logger)

	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Second)
}

func (s *InitChallengeSuite) TearDownTest() {
	s.cancel()
}

func (s *InitChallengeSuite) TestInitChallenge() {
	id := uuid.New()
	challengeInfo := dto.NewChallengeInfo(2, 10)
	s.repo.On("SetChallengeInfo", id, challengeInfo).Return(nil)

	err := s.uc.Init(s.ctx, id, challengeInfo)

	s.Require().NoError(err)
}

func (s *InitChallengeSuite) TestInitChallenge_error() {
	id := uuid.New()
	challengeInfo := dto.NewChallengeInfo(2, 10)
	s.repo.On("SetChallengeInfo", id, challengeInfo).Return(errors.New("db connection error"))

	err := s.uc.Init(s.ctx, id, challengeInfo)

	s.Require().Error(err)
	s.EqualError(err, "failed to set challenge info: db connection error")
}
