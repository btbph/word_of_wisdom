package usecase

import (
	"context"
	"fmt"
	"github.com/btbph/word_of_wisdom/internal/dto"
	"github.com/google/uuid"
	"log/slog"
)

type InitChallengeRepo interface {
	SetChallengeInfo(ctx context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error
}

type InitChallenge struct {
	repo   InitChallengeRepo
	logger *slog.Logger
}

func NewInitChallenge(repo InitChallengeRepo, logger *slog.Logger) *InitChallenge {
	return &InitChallenge{
		repo:   repo,
		logger: logger,
	}
}

func (uc *InitChallenge) Init(ctx context.Context, ID uuid.UUID, challengeInfo dto.ChallengeInfo) error {
	if err := uc.repo.SetChallengeInfo(ctx, ID, challengeInfo); err != nil {
		uc.logger.Error("failed to set challenge info", "error", err)
		return fmt.Errorf("failed to set challenge info: %w", err)
	}

	return nil
}
