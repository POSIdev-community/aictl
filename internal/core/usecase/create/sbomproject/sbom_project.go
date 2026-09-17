package sbomproject

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectByName(ctx context.Context, projectName string) (*project.Project, error)
	CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error)
	UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error
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

func (u *UseCase) Execute(ctx context.Context, projectName, sbomPath string, safe bool) error {
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
		if existing.Type.IsSource() {
			return project.ErrNameUsedBySource
		}

		if !safe {
			return project.ErrProjectNameUsedBySbom
		}

		if err = u.aiAdapter.UpdateSbom(ctx, existing.Id, sbomPath); err != nil {
			return fmt.Errorf("update sbom: %w", err)
		}

		u.cliAdapter.ShowTextf(ctx, "SBOM updated, id '%v'", existing.Id.String())
		u.cliAdapter.ReturnText(ctx, existing.Id.String())

		return nil
	}

	projectId, err := u.aiAdapter.CreateSbomProject(ctx, projectName)
	if err != nil {
		return fmt.Errorf("create sbom project: %w", err)
	}

	if err = u.aiAdapter.UpdateSbom(ctx, *projectId, sbomPath); err != nil {
		return fmt.Errorf("update sbom: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "project '%v' created, id '%v'", projectName, projectId.String())
	u.cliAdapter.ReturnText(ctx, projectId.String())

	return nil
}
