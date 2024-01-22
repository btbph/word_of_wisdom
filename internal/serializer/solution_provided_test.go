package serializer

import (
	"bytes"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"testing"
)

type SolutionProvidedSuite struct {
	suite.Suite

	ser *SolutionProvided

	logger *slog.Logger
}

func TestSolutionProvidedSuite(t *testing.T) {
	suite.Run(t, &SolutionProvidedSuite{})
}

func (s *SolutionProvidedSuite) SetupSuite() {
	s.logger = slog.Default()
	s.ser = NewSolutionProvided(s.logger)
}

func (s *SolutionProvidedSuite) TestSerializationSolutionProvided() {
	req := request.NewSolutionProvided("solution")

	result, err := s.ser.Serialize(convertToReader(req))

	s.Require().NoError(err)
	s.Equal(req, result)
}

func (s *SolutionProvidedSuite) TestSerializationSolutionProvided_wrongType() {
	req := request.NewSolutionProvided("solution")
	req.Type = dto.RequestChallenge

	result, err := s.ser.Serialize(convertToReader(req))

	s.Require().Error(err)
	s.Require().Empty(result)
	s.EqualError(err, "wrong request type")
}

func (s *SolutionProvidedSuite) TestSerializationSolutionProvided_wrongRequest() {
	req := bytes.NewBuffer([]byte{1, 2, 3, 4})

	result, err := s.ser.Serialize(req)

	s.Require().Error(err)
	s.Require().Empty(result)
	s.EqualError(err, "wrong request type")
}
