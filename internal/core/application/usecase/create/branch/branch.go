package branch

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error)
	CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTarget string) (*uuid.UUID, error)
}

type CLI interface {
	ReturnText(text string)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, branchName, scanTarget string, safe bool) error {
	if safe {

		branches, err := u.aiAdapter.GetBranches(ctx, cfg.ProjectId())
		if err != nil {
			return fmt.Errorf("get branches useCase error: %w", err)
		}

		for _, b := range branches {
			if b.Name == branchName {
				u.cliAdapter.ReturnText(b.Id.String())
				return nil
			}
		}
	}

	branchId, err := u.aiAdapter.CreateBranch(ctx, cfg.ProjectId(), branchName, scanTarget)
	if err != nil {
		return fmt.Errorf("usecase create branch: %w", err)
	}

	u.cliAdapter.ReturnText(branchId.String())

	return nil
}
