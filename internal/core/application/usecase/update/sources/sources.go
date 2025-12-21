package sources

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	Initialize(ctx context.Context) error
	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, sourcePath string) error
}

type CLI interface {
	ShowText(ctx context.Context, text string)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, sourcePath string) error {
	err := u.aiAdapter.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize with jwt retry: %w", err)
	}

	u.cliAdapter.ShowText(ctx, "start updating sources")

	err = u.aiAdapter.UpdateSources(ctx, u.cfg.ProjectId(), u.cfg.BranchId(), sourcePath)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowText(ctx, "sources updated")

	return nil
}
