package dto

type ChallengeInfo struct {
	ZeroBits   int
	SaltLength int
}

func NewChallengeInfo(
	zeroBits int,
	saltLength int,
) ChallengeInfo {
	return ChallengeInfo{
		ZeroBits:   zeroBits,
		SaltLength: saltLength,
	}
}
