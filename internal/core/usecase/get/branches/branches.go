package branches

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error)
}

type CLI interface {
	ShowBranches(ctx context.Context, branches []branch.Branch)
	ShowBranchesQuite(ctx context.Context, branches []branch.Branch)
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

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotGetBranchesOnSbom); err != nil {
		return err
	}

	branches, err := u.aiAdapter.GetBranches(ctx, u.cfg.ProjectId())
	if err != nil {
		return fmt.Errorf("get branches: %w", err)
	}

	filteredBranches := make([]branch.Branch, 0, len(branches))
	if filter.Empty() {
		filteredBranches = branches
	} else {
		for _, b := range branches {
			if filter.Execute(b.Name) {
				filteredBranches = append(filteredBranches, b)
			}
		}
	}

	if quite {
		u.cliAdapter.ShowBranchesQuite(ctx, filteredBranches)
	} else {
		u.cliAdapter.ShowBranches(ctx, filteredBranches)
	}

	return nil
}
