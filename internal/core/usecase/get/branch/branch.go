package branch

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetBranch(ctx context.Context, branchId uuid.UUID) (*branch.Branch, error)
}

type CLI interface {
	ShowBranches(ctx context.Context, branches []branch.Branch)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, branchId uuid.UUID) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	// When project is in context, reject SBOM projects (branch API has no project id).
	if u.cfg != nil && u.cfg.ProjectId() != uuid.Nil {
		if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotGetBranchOnSbom); err != nil {
			return err
		}
	}

	b, err := u.aiAdapter.GetBranch(ctx, branchId)
	if err != nil {
		return fmt.Errorf("get branch: %w", err)
	}

	u.cliAdapter.ShowBranches(ctx, []branch.Branch{*b})

	return nil
}
