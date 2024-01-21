package hashcash

import (
	"crypto/sha256"
	"github.com/btbph/word_of_wisdom/internal/hashcash/mock"
	"github.com/stretchr/testify/suite"
	"hash"
	"math"
	"testing"
	"time"
)

const (
	testZeroBits   = 8
	testSaltLength = 8
)

var n = time.Date(2024, 1, 20, 16, 00, 00, 00, time.UTC)

type HashcashSuite struct {
	suite.Suite

	hashcash *Hashcash

	clock *mock.Clock

	hasher hash.Hash
}

func TestHashcashSuite(t *testing.T) {
	suite.Run(t, &HashcashSuite{})
}

func (s *HashcashSuite) SetupSuite() {
	s.hasher = sha256.New()
}

func (s *HashcashSuite) SetupTest() {
	s.clock = &mock.Clock{}
	s.hashcash = NewHashcash(testZeroBits, testSaltLength, s.clock, s.hasher)
}

func (s *HashcashSuite) TestGenerate() {
	resource := "resource"
	s.clock.On("Now").Return(n)

	result := s.hashcash.Generate(resource)

	s.assertValidHashcashString(result)
}

func (s *HashcashSuite) TestGenerate_emptyResource() {
	resource := ""
	s.clock.On("Now").Return(n)

	result := s.hashcash.Generate(resource)

	s.assertValidHashcashString(result)
}

func (s *HashcashSuite) assertValidHashcashString(hashcashString string) {
	s.hasher.Reset()
	s.hasher.Write([]byte(hashcashString))
	generatedHash := s.hasher.Sum(nil)
	hexDigits := int(math.Ceil(float64(testZeroBits) / 4))

	i := 0
	for i < hexDigits {
		s.Require().Zero(generatedHash[i])
		i++
	}
}
