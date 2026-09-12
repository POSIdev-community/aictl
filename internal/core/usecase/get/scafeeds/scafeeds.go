package scafeeds

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScaFeeds(ctx context.Context, statuses []scafeeds.Status) ([]scafeeds.Package, error)
}

type CLI interface {
	ShowScaFeeds(ctx context.Context, packages []scafeeds.Package)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}
	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, statuses []scafeeds.Status) error {
	if err := u.aiAdapter.InitializeWithRetry(ctx); err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	pkgs, err := u.aiAdapter.GetScaFeeds(ctx, statuses)
	if err != nil {
		return fmt.Errorf("get sca feeds: %w", err)
	}

	u.cliAdapter.ShowScaFeeds(ctx, pkgs)

	return nil
}
