package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type CheckSolutionRepo interface {
	RemoveChallengeInfo(ctx context.Context, ID uuid.UUID) error
	CheckSolutionPresence(ctx context.Context, solution string) (bool, error)
	SaveSolution(ctx context.Context, solution string) error
	RandomQuote(ctx context.Context) (string, error)
}

type CheckSolutionValidator interface {
	Validate(solution string) error
}

type CheckSolution struct {
	repo      CheckSolutionRepo
	validator CheckSolutionValidator
	logger    *slog.Logger
}

func NewCheckSolution(
	repo CheckSolutionRepo,
	validator CheckSolutionValidator,
	logger *slog.Logger,
) *CheckSolution {
	return &CheckSolution{
		repo:      repo,
		validator: validator,
		logger:    logger,
	}
}

func (uc *CheckSolution) Check(ctx context.Context, ID uuid.UUID, solution string) (string, error) {
	if err := uc.validator.Validate(solution); err != nil {
		uc.logger.Warn("solution validation failed", "error", err)
		return "", fmt.Errorf("solution validation failed: %w", err)
	}

	present, err := uc.repo.CheckSolutionPresence(ctx, solution)
	if err != nil {
		uc.logger.Error("failed to check solution presence", "error", err)
		return "", fmt.Errorf("failed to check solution presence: %w", err)
	}

	if present {
		uc.logger.Warn("current solution already presents", "solution", solution)
		return "", errors.New("current solution already presents")
	}

	if err = uc.repo.SaveSolution(ctx, solution); err != nil {
		uc.logger.Error("failed to save solution", "error", err)
		return "", fmt.Errorf("failed to save solution: %w", err)
	}

	quote, err := uc.repo.RandomQuote(ctx)
	if err != nil {
		uc.logger.Error("failed to get random quote", "error", err)
		return "", fmt.Errorf("failed to get random quote: %w", err)
	}

	go func() {
		if err = uc.repo.RemoveChallengeInfo(ctx, ID); err != nil {
			uc.logger.Error("failed to remove challenge info", "error", err)
		}
	}()

	return quote, nil
}
