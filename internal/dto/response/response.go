package response

import (
	"github.com/btbph/word_of_wisdom/internal/dto"
)

type RequestChallenge struct {
	Type       dto.Type `json:"action"`
	ZeroBits   int      `json:"zeroBits"`
	SaltLength int      `json:"saltLength"`
}

func NewRequestChallengeResponse(zeroBits, saltLength int) RequestChallenge {
	return RequestChallenge{
		Type:       dto.ReturnChalange,
		ZeroBits:   zeroBits,
		SaltLength: saltLength,
	}
}

type SolutionProvided struct {
	Type  dto.Type `json:"type"`
	Quote string   `json:"quote"`
}

func NewSolutionProvided(quote string) *SolutionProvided {
	return &SolutionProvided{
		Type:  dto.QuoteProvided,
		Quote: quote,
	}
}
