package languages

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
	UpdateProjectLanguages(ctx context.Context, projectId uuid.UUID) error
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

func (u *UseCase) Execute(ctx context.Context) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotGetSetSettingsOnSbom); err != nil {
		return err
	}

	u.cliAdapter.ShowText(ctx, "updating project languages from uploaded sources")

	if err := u.aiAdapter.UpdateProjectLanguages(ctx, u.cfg.ProjectId()); err != nil {
		return fmt.Errorf("update project languages: %w", err)
	}

	u.cliAdapter.ShowText(ctx, "project languages updated")

	return nil
}
