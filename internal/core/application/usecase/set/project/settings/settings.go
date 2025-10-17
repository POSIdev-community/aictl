package settings

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	GetDefaultSettings(ctx context.Context) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
}

type CLI interface {
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, projectID uuid.UUID, aiProj *aiproj.AIProj) error {
	scanSettings, err := u.aiAdapter.GetDefaultSettings(ctx)
	if err != nil {
		return fmt.Errorf("get default settings: %w", err)
	}

	scanSettings.UpdateFromAIProj(aiProj)

	if err := u.aiAdapter.SetProjectSettings(ctx, projectID, &scanSettings); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	return nil
}
