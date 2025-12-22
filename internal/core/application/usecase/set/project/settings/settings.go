package settings

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
}

type CLI interface {
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

func (u *UseCase) Execute(ctx context.Context, aiProj *aiproj.AIProj) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	scanSettings, err := u.aiAdapter.GetDefaultSettings(ctx)
	if err != nil {
		return fmt.Errorf("get default settings: %w", err)
	}

	scanSettings.UpdateFromAIProj(aiProj)

	if err := u.aiAdapter.SetProjectSettings(ctx, u.cfg.ProjectId(), &scanSettings); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	return nil
}
