package branch

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetBranch(ctx context.Context, branchId uuid.UUID) (*branch.Branch, error)
}

type CLI interface {
	ShowBranches(ctx context.Context, branches []branch.Branch)
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, branchId uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	b, err := u.aiAdapter.GetBranch(ctx, branchId)
	if err != nil {
		return fmt.Errorf("get branch: %w", err)
	}

	u.cliAdapter.ShowBranches(ctx, []branch.Branch{*b})

	return nil
}
