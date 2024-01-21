package hashcash

import (
	"crypto/sha256"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/hashcash/mock"
	"github.com/stretchr/testify/suite"
	"hash"
	"testing"
	"time"
)

type ValidatorSuite struct {
	suite.Suite

	validator Validator
	clock     *mock.Clock
	hasher    hash.Hash

	expiryTime       time.Duration
	expectedResource string
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, &ValidatorSuite{})
}

func (s *ValidatorSuite) SetupSuite() {
	s.hasher = sha256.New()
	s.expiryTime = 24 * time.Hour
	s.expectedResource = "resource"
}

func (s *ValidatorSuite) SetupTest() {
	s.clock = &mock.Clock{}
	challengeInfo := dto.ChallengeInfo{
		ZeroBits:   testZeroBits,
		SaltLength: testSaltLength,
	}
	s.validator = NewValidator(challengeInfo, s.clock, s.expiryTime, s.hasher, s.expectedResource)
}

func (s *ValidatorSuite) TestValidate() {
	s.clock.On("Now").Return(n)
	resource := "1:8:240120:resource::i29+cKjq:7bde"

	err := s.validator.Validate(resource)

	s.Require().NoError(err)
}

func (s *ValidatorSuite) TestValidate_errors() {
	tt := []struct {
		name  string
		setup func() string
		err   string
	}{
		{
			name: "wrong format",
			setup: func() string {
				return "1:8:240120:resource::i29+cKjq"
			},
			err: "wrong format of provided string",
		},
		{
			name: "wrong version",
			setup: func() string {
				return "2:8:240120:resource::i29+cKjq:7bde"
			},
			err: "expected first version",
		},
		{
			name: "wrong bit size",
			setup: func() string {
				return "1:9:240120:resource::i29+cKjq:7bde"
			},
			err: "bit size doesn't match",
		},
		{
			name: "wrong date format",
			setup: func() string {
				return "1:8:20240120:resource::i29+cKjq:7bde"
			},
			err: "wrong date format",
		},
		{
			name: "future date",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:240122:resource::i29+cKjq:7bde"
			},
			err: "future dates aren't allowed",
		},
		{
			name: "expired date",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:230120:resource::i29+cKjq:7bde"
			},
			err: "provided string expired",
		},
		{
			name: "wrong resource",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:240120:notResource::i29+cKjq:7bde"
			},
			err: "string doesn't contain needed resource",
		},
		{
			name: "extension presents",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:240120:resource:1:i29+cKjq:7bde"
			},
			err: "extenstion should be ignored",
		},
		{
			name: "wrong salt length",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:240120:resource::i29+cKjqa:7bde"
			},
			err: "salt length doesn't match",
		},
		{
			name: "leading bits aren't zeros",
			setup: func() string {
				s.clock.On("Now").Return(n)
				return "1:8:240120:resource::VEFhgq6g:0"
			},
			err: "leading zeros aren't zeroes",
		},
	}

	for _, tc := range tt {
		s.Run(tc.name, func() {
			s.SetupTest()
			resource := tc.setup()

			err := s.validator.Validate(resource)

			s.Require().Error(err)
			s.EqualError(err, tc.err)
		})
	}
}
