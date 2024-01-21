package hashcash

import (
	"errors"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"hash"
	"strconv"
	"strings"
	"time"
)

type Validator struct {
	challengeInfo    dto.ChallengeInfo
	clock            Clock
	expiryTime       time.Duration
	hasher           hash.Hash
	expectedResource string
}

func NewValidator(
	challengeInfo dto.ChallengeInfo,
	clock Clock,
	expiryTime time.Duration,
	hasher hash.Hash,
	expectedResource string,
) Validator {
	return Validator{
		challengeInfo:    challengeInfo,
		clock:            clock,
		expiryTime:       expiryTime,
		hasher:           hasher,
		expectedResource: expectedResource,
	}
}

func (v Validator) Validate(str string) error {
	parts := strings.Split(str, ":")
	if len(parts) != 7 {
		return errors.New("wrong format of provided string")
	}

	if parts[0] != "1" {
		return errors.New("expected first version")
	}

	if parts[1] != strconv.Itoa(v.challengeInfo.ZeroBits) {
		return errors.New("bit size doesn't match")
	}

	t, err := time.Parse("060102", parts[2])
	if err != nil {
		return errors.New("wrong date format")
	}

	now := v.clock.Now()
	if t.After(now) {
		return errors.New("future dates aren't allowed")
	}

	if now.Sub(t) > v.expiryTime {
		return errors.New("provided string expired")
	}

	if parts[3] != v.expectedResource {
		return errors.New("string doesn't contain needed resource")
	}

	if parts[4] != "" {
		return errors.New("extenstion should be ignored")
	}

	if len(parts[5]) != v.challengeInfo.SaltLength {
		return errors.New("salt length doesn't match")
	}

	if !checkZeroes(v.hasher, v.challengeInfo.ZeroBits, str) {
		return errors.New("leading zeros aren't zeroes")
	}

	return nil
}
