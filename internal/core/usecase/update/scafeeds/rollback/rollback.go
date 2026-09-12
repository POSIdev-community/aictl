package rollback

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	RollbackScaFeeds(ctx context.Context) (scafeeds.Package, error)
}

type CLI interface {
	AskConfirmation(ctx context.Context, question string) (bool, error)
	ShowText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
	ReturnText(ctx context.Context, text string)
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

func (u *UseCase) Execute(ctx context.Context, skipConfirm bool) error {
	if err := u.aiAdapter.InitializeWithRetry(ctx); err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	ok := skipConfirm
	if !ok {
		var err error
		ok, err = u.cliAdapter.AskConfirmation(ctx, "Roll back current SCA feeds package?")
		if err != nil {
			return err
		}
	}

	if !ok {
		u.cliAdapter.ShowText(ctx, "Cancelled")

		return nil
	}

	pkg, err := u.aiAdapter.RollbackScaFeeds(ctx)
	if err != nil {
		return fmt.Errorf("rollback sca feeds: %w", err)
	}

	u.cliAdapter.ReturnText(ctx, pkg.Version)
	u.cliAdapter.ShowTextf(ctx, "SCA feeds rolled back to version '%s'", pkg.Version)

	return nil
}
