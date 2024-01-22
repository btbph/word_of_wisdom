package request

import (
	"encoding/json"
	"github.com/btbph/word_of_wisdom/internal/dto"
)

type RequestChallenge struct {
	Type dto.Type `json:"action"`
}

func NewRequestChallenge() RequestChallenge {
	return RequestChallenge{
		Type: dto.RequestChallenge,
	}
}

func MarshalRequestChallenge(req RequestChallenge) ([]byte, error) {
	return json.Marshal(req)
}

type SolutionProvided struct {
	Type     dto.Type `json:"action"`
	Solution string   `json:"solution"`
}

func NewSolutionProvided(solution string) SolutionProvided {
	return SolutionProvided{
		Type:     dto.SolutionProvided,
		Solution: solution,
	}
}

func MarshalSolutionProvided(req SolutionProvided) ([]byte, error) {
	return json.Marshal(req)
}
