package serializer

import (
	"bytes"
	"encoding/json"
	"github.com/btbph/word_of_wisdom/internal/dto/request"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"testing"
)

type RequestChallengeSuite struct {
	suite.Suite

	ser *RequestChallenge

	logger *slog.Logger
}

func TestRequestChallengeSuite(t *testing.T) {
	suite.Run(t, &RequestChallengeSuite{})
}

func (s *RequestChallengeSuite) SetupSuite() {
	s.logger = slog.Default()
	s.ser = NewRequestChallenge(s.logger)
}

func (s *RequestChallengeSuite) TestSerializationRequestChallenge() {
	req := request.NewRequestChallenge()

	result, err := s.ser.Serialize(convertToReader(req))

	s.Require().NoError(err)
	s.Equal(req, result)
}

func (s *RequestChallengeSuite) TestSerializationRequestChallenge_wrongType() {
	req := request.NewRequestChallenge()
	req.Type = 10

	result, err := s.ser.Serialize(convertToReader(req))

	s.Require().Error(err)
	s.Require().Empty(result)
	s.EqualError(err, "expect request challenge")
}

func (s *RequestChallengeSuite) TestSerializationRequestChallenge_wrongRequest() {
	req := bytes.NewBuffer([]byte{1, 2, 3, 4})

	result, err := s.ser.Serialize(req)

	s.Require().Error(err)
	s.Require().Empty(result)
	s.EqualError(err, "expect request challenge")
}

func convertToReader[T any](req T) io.Reader {
	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}
