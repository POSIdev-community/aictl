package project

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error)
	CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error)
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
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, projectName string, safe bool) error {
	var (
		projectId *uuid.UUID
		err       error
	)

	err = u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize with jwt retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "creating project '%v'", projectName)

	if safe {
		projectId, err = u.aiAdapter.GetProjectId(ctx, projectName)
		if err != nil {
			return err
		}
	}

	if projectId != nil {
		u.cliAdapter.ShowTextf(ctx, "project '%v' already exixts, id '%v'", projectName, projectId.String())
		u.cliAdapter.ReturnText(ctx, projectId.String())
		return nil
	}

	projectId, err = u.aiAdapter.CreateProject(ctx, projectName)
	if err != nil {
		return fmt.Errorf("usecase create branch: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "project '%v' created, id '%v'", projectName, projectId.String())
	u.cliAdapter.ReturnText(ctx, projectId.String())

	return nil
}
