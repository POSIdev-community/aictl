package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	domainproject "github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetLicense(ctx context.Context) (*domainlicense.License, error)
	GetProject(ctx context.Context, projectId uuid.UUID) (*domainproject.Project, error)
	GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error)
	SetProjectSettings(ctx context.Context, projectId uuid.UUID, settings *settings.ScanSettings) error
	StartScanProject(ctx context.Context, projectId uuid.UUID, scanLabel string, scanType scantype.Type) (uuid.UUID, error)
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

func (u *UseCase) Execute(ctx context.Context, scanLabel string, scanType scantype.Type) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err = usecaseutils.RequireSourceProject(ctx, u.aiAdapter, u.cfg.ProjectId(), domainproject.ErrCannotRunProjectScanOnSbom); err != nil {
		return err
	}

	if err = usecaseutils.ApplyLicenseToProjectSettings(ctx, u.aiAdapter, u.cliAdapter, u.cfg.ProjectId(), usecaseutils.ApplyLicenseOptions{}); err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "starting scan, project-id '%v', branch-id '%v'", u.cfg.ProjectId(), u.cfg.BranchId())

	scanResultId, err := u.aiAdapter.StartScanProject(ctx, u.cfg.ProjectId(), scanLabel, scanType)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "scan started, scan-id '%v'", scanResultId)
	u.cliAdapter.ReturnText(ctx, scanResultId.String())

	return nil
}
