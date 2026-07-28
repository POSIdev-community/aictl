package sources

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type AI interface {
	Initialize(ctx context.Context) error
	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, sourcePath string, exclusions gitignore.Exclusions, tempDir string) error
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
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, sourcePath string, exclusions gitignore.Exclusions, tempDir string) error {
	err := u.aiAdapter.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	u.cliAdapter.ShowText(ctx, "start updating sources")

	err = u.aiAdapter.UpdateSources(ctx, u.cfg.ProjectId(), u.cfg.BranchId(), sourcePath, exclusions, tempDir)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowText(ctx, "sources updated")

	return nil
}
