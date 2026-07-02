package settings

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetVersion(ctx context.Context) (version.Version, error)
	GetDefaultSettings(ctx context.Context) (domainsettings.ScanSettings, error)
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

func (u *UseCase) Execute(ctx context.Context, rawAiproj []byte) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	serverVersion, err := u.aiAdapter.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("get server version: %w", err)
	}

	aiProj, err := aiproj.ParseAndMigrateForServer(rawAiproj, serverVersion)
	if err != nil {
		return fmt.Errorf("parse and migrate aiproj: %w", err)
	}

	scanSettings, err := u.aiAdapter.GetDefaultSettings(ctx)
	if err != nil {
		return fmt.Errorf("get default settings: %w", err)
	}

	err = scanSettings.UpdateFromParsed(aiProj)
	if err != nil {
		return fmt.Errorf("update scan settings from aiproj: %w", err)
	}

	if err := u.aiAdapter.SetProjectSettings(ctx, u.cfg.ProjectId(), &scanSettings); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	return nil
}
