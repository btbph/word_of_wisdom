package dto

type Type int

const (
	RequestChallenge Type = iota + 1
	ReturnChalange
	SolutionProvided
	QuoteProvided
)
