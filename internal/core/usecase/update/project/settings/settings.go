package settings

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetVersion(ctx context.Context) (version.Version, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (domainsettings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *domainsettings.ScanSettings) error
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
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, patch domainsettings.ProjectSettingsPatch) error {
	if patch.Priority == nil && patch.PreferredAgents == nil && patch.PreferredAgentsOnly == nil {
		return validation.NewError("at least one of --priority, --agents, --preferred-agents-only or --no-preferred-agents-only is required")
	}

	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotGetSetSettingsOnSbom); err != nil {
		return err
	}

	serverVersion, err := u.aiAdapter.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("get server version: %w", err)
	}

	if err := usecaseutils.RequireProjectScanSettings(serverVersion); err != nil {
		return err
	}

	current, err := u.aiAdapter.GetProjectSettings(ctx, u.cfg.ProjectId())
	if err != nil {
		return fmt.Errorf("get project settings: %w", err)
	}

	current.Patch(patch)

	if err := u.aiAdapter.SetProjectSettings(ctx, u.cfg.ProjectId(), &current); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	return nil
}
