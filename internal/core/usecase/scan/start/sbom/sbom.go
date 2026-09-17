package sbom

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetLicense(ctx context.Context) (*domainlicense.License, error)
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
	StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error)
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

func (u *UseCase) Execute(ctx context.Context, scanLabel string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSbomProject(ctx, u.aiAdapter, u.cfg.ProjectId(), project.ErrCannotRunSbomScanOnSource); err != nil {
		return err
	}

	if err = usecaseutils.ApplyLicenseToProjectSettings(ctx, u.aiAdapter, u.cliAdapter, u.cfg.ProjectId(), usecaseutils.ApplyLicenseOptions{
		SkipLanguageCheck: true,
	}); err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "starting scan, project-id '%v'", u.cfg.ProjectId())

	scanResultId, err := u.aiAdapter.StartScanSbom(ctx, u.cfg.ProjectId(), scanLabel)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "scan started, scan-id '%v'", scanResultId)
	u.cliAdapter.ReturnText(ctx, scanResultId.String())

	return nil
}
