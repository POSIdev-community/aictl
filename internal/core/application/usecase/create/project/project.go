package project

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	GetProjectId(ctx context.Context, projectName string) (*uuid.UUID, error)
	CreateProject(ctx context.Context, projectName string) (*uuid.UUID, error)
}

type CLI interface {
	ReturnText(text string)
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

	if safe {
		projectId, err = u.aiAdapter.GetProjectId(ctx, projectName)
		if err != nil {
			return err
		}
	}

	if projectId == nil {
		projectId, err = u.aiAdapter.CreateProject(ctx, projectName)
		if err != nil {
			return fmt.Errorf("usecase create branch: %w", err)
		}
	}

	u.cliAdapter.ReturnText(projectId.String())

	return nil
}
