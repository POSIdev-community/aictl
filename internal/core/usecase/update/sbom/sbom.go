package sbom

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error
}

type CLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
	ReturnText(ctx context.Context, text string)
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

func (u *UseCase) Execute(ctx context.Context, sbomPath string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSbomProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotUpdateSbomOnSource); err != nil {
		return err
	}

	if err = u.aiAdapter.UpdateSbom(ctx, u.cfg.ProjectId(), sbomPath); err != nil {
		return fmt.Errorf("update sbom: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "SBOM updated, id '%v'", u.cfg.ProjectId())
	u.cliAdapter.ReturnText(ctx, u.cfg.ProjectId().String())

	return nil
}
