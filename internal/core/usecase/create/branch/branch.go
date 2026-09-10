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
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type AI interface {
	Initialize(ctx context.Context) error
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetBranches(ctx context.Context, projectId uuid.UUID) ([]branch.Branch, error)
	CreateBranch(ctx context.Context, projectId uuid.UUID, branchName, scanTarget string, exclusions gitignore.Exclusions, tempDir string) (*uuid.UUID, error)
}

type CLI interface {
	ReturnText(ctx context.Context, text string)
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

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, branchName, scanTarget string, safe bool, exclusions gitignore.Exclusions, tempDir string) error {
	err := u.aiAdapter.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, cfg.ProjectId(), project.ErrCannotCreateBranchOnSbom); err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "creating branch '%v'", branchName)

	if safe {

		branches, err := u.aiAdapter.GetBranches(ctx, cfg.ProjectId())
		if err != nil {
			return fmt.Errorf("get branches: %w", err)
		}

		for _, b := range branches {
			if b.Name == branchName {
				u.cliAdapter.ShowTextf(ctx, "branch '%v' already exists, id '%v'", branchName, b.Id.String())
				u.cliAdapter.ReturnText(ctx, b.Id.String())
				return nil
			}
		}
	}

	branchId, err := u.aiAdapter.CreateBranch(ctx, cfg.ProjectId(), branchName, scanTarget, exclusions, tempDir)
	if err != nil {
		return fmt.Errorf("create branch: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "branch '%v' created, id '%v'", branchName, branchId.String())
	u.cliAdapter.ReturnText(ctx, branchId.String())

	return nil
}
