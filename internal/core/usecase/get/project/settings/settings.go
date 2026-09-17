package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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
}

type CLI interface {
	ShowProjectSettings(ctx context.Context, view domainsettings.ProjectSettingsView)
	ShowReader(io.Reader) error
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

func (u *UseCase) Execute(ctx context.Context, jsonOutput bool) error {
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

	view := current.ProjectSettingsView()

	if !jsonOutput {
		u.cliAdapter.ShowProjectSettings(ctx, view)

		return nil
	}

	jsonData, err := json.MarshalIndent(view, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal project settings: %w", err)
	}
	jsonData = append(jsonData, '\n')

	reader := bytes.NewReader(jsonData)
	if err := u.cliAdapter.ShowReader(reader); err != nil {
		return fmt.Errorf("print project settings: %w", err)
	}

	return nil
}
