package scafeeds

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	UpdateScaFeeds(ctx context.Context, path, version string) error
}

type CLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
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

func (u *UseCase) Execute(ctx context.Context, path, version string) error {
	if err := u.aiAdapter.InitializeWithRetry(ctx); err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err := u.aiAdapter.UpdateScaFeeds(ctx, path, version); err != nil {
		return fmt.Errorf("update sca feeds: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "SCA feeds updated, version '%s'", version)

	return nil
}
