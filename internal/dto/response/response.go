package response

import (
	"encoding/json"
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

func MarshalRequestChallenge(res RequestChallenge) ([]byte, error) {
	return json.Marshal(res)
}

type SolutionProvided struct {
	Type  dto.Type `json:"type"`
	Quote string   `json:"quote"`
}

func NewSolutionProvided(quote string) SolutionProvided {
	return SolutionProvided{
		Type:  dto.QuoteProvided,
		Quote: quote,
	}
}

func MarshalSolutionProvided(res SolutionProvided) ([]byte, error) {
	return json.Marshal(res)
}
