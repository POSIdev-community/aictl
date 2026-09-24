package policycheck

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
}

type CLI interface {
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

func (u *UseCase) Execute(ctx context.Context) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	model, err := usecaseutils.LoadProjectPolicies(ctx, u.aiAdapter, u.cfg.ProjectId())
	if err != nil {
		return err
	}

	u.cliAdapter.ReturnText(ctx, strconv.FormatBool(model.CheckAccordance))

	return nil
}
