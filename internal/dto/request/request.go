package request

import "github.com/btbph/word_of_wisdom/internal/dto"

type RequestChallenge struct {
	Type dto.Type `json:"action"`
}

func NewRequestChallenge() RequestChallenge {
	return RequestChallenge{
		Type: dto.RequestChallenge,
	}
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
