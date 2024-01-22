package usecase

import (
	"context"
	"errors"
	"github.com/btbph/word_of_wisdom/internal/usecase/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"testing"
	"time"
)

type CheckSolutionSuite struct {
	suite.Suite

	uc        *CheckSolution
	repo      *mock.Repo
	validator *mock.Validator

	logger *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func TestCheckSolutionSuite(t *testing.T) {
	suite.Run(t, &CheckSolutionSuite{})
}

func (s *CheckSolutionSuite) SetupSuite() {
	s.logger = slog.Default()
}

func (s *CheckSolutionSuite) SetupTest() {
	s.repo = &mock.Repo{}
	s.validator = &mock.Validator{}
	s.uc = NewCheckSolution(s.repo, s.validator, s.logger)

	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Second)
}

func (s *CheckSolutionSuite) TearDownTest() {
	s.cancel()
}

func (s *CheckSolutionSuite) TestCheckSolution() {
	id := uuid.New()
	const (
		solution = "solution"
		quote    = "random quote"
	)
	s.validator.On("Validate", solution).Return(nil)
	s.repo.On("CheckSolutionPresence", solution).Return(false, nil)
	s.repo.On("SaveSolution", solution).Return(nil)
	s.repo.On("RandomQuote").Return(quote, nil)
	s.repo.On("RemoveChallengeInfo", id).Return(nil)

	result, err := s.uc.Check(s.ctx, id, solution)

	s.Require().NoError(err)
	s.Equal(result, quote)
	s.EventuallyWithT(func(collect *assert.CollectT) {
		s.repo.AssertNumberOfCalls(s.T(), "RemoveChallengeInfo", 1)
	}, time.Second, time.Millisecond)
}

func (s *CheckSolutionSuite) TestCheckSolution_error() {
	tt := []struct {
		name  string
		setup func() string
		error string
	}{
		{
			name: "failed to validate solution",
			setup: func() string {
				const solution = "solution"
				s.validator.On("Validate", solution).Return(errors.New("validation failed"))
				return solution
			},
			error: "solution validation failed: validation failed",
		},
		{
			name: "failed to check solution presense",
			setup: func() string {
				const solution = "solution"
				s.validator.On("Validate", solution).Return(nil)
				s.repo.On("CheckSolutionPresence", solution).Return(false, errors.New("db error"))
				return solution
			},
			error: "failed to check solution presence: db error",
		},
		{
			name: "solution present",
			setup: func() string {
				const solution = "solution"
				s.validator.On("Validate", solution).Return(nil)
				s.repo.On("CheckSolutionPresence", solution).Return(true, nil)
				return solution
			},
			error: "current solution already presents",
		},
		{
			name: "failed to save solution",
			setup: func() string {
				const solution = "solution"
				s.validator.On("Validate", solution).Return(nil)
				s.repo.On("CheckSolutionPresence", solution).Return(false, nil)
				s.repo.On("SaveSolution", solution).Return(errors.New("db error"))
				return solution
			},
			error: "failed to save solution: db error",
		},
		{
			name: "failed to get random quote",
			setup: func() string {
				const solution = "solution"
				s.validator.On("Validate", solution).Return(nil)
				s.repo.On("CheckSolutionPresence", solution).Return(false, nil)
				s.repo.On("SaveSolution", solution).Return(nil)
				s.repo.On("RandomQuote").Return("", errors.New("db error"))
				return solution
			},
			error: "failed to get random quote: db error",
		},
	}

	for _, tc := range tt {
		s.Run(tc.name, func() {
			s.SetupTest()
			solution := tc.setup()

			result, err := s.uc.Check(s.ctx, uuid.New(), solution)

			s.Require().Error(err)
			s.Require().Empty(result)
			s.EqualError(err, tc.error)
			s.TearDownTest()
		})
	}
}
