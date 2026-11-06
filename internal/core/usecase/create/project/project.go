package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
	domainproject "github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectByName(ctx context.Context, projectName string) (*domainproject.Project, error)
	GetLicense(ctx context.Context) (*domainlicense.License, error)
	CreateProject(ctx context.Context, projectName string, language settings.Language) (*uuid.UUID, error)
}

type CLI interface {
	ReturnText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, projectName string, safe bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "creating project '%v'", projectName)

	existing, err := u.aiAdapter.GetProjectByName(ctx, projectName)
	if err != nil {
		return fmt.Errorf("get project by name: %w", err)
	}

	if existing != nil {
		if existing.Type.IsSbom() {
			return domainproject.ErrNameUsedBySbom
		}

		if !safe {
			return domainproject.ErrProjectNameUsedBySource
		}

		u.cliAdapter.ShowTextf(ctx, "project '%v' already exists, id '%v'", projectName, existing.Id.String())
		u.cliAdapter.ReturnText(ctx, existing.Id.String())

		return nil
	}

	lic, err := u.aiAdapter.GetLicense(ctx)
	if err != nil {
		return fmt.Errorf("get license: %w", err)
	}

	language, err := lic.ProjectLanguage()
	if err != nil {
		return err
	}

	projectId, err := u.aiAdapter.CreateProject(ctx, projectName, language)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "project '%v' created, id '%v'", projectName, projectId.String())
	u.cliAdapter.ReturnText(ctx, projectId.String())

	return nil
}
