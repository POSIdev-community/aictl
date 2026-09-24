package policycheck

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	usecaseutils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
	SetProjectPolicies(ctx context.Context, projectId uuid.UUID, rawJSON []byte) error
}

type CLI interface{}

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

func (u *UseCase) Execute(ctx context.Context, enabled bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	if err := usecaseutils.SetProjectPolicyCheck(ctx, u.aiAdapter, u.cfg.ProjectId(), enabled); err != nil {
		return fmt.Errorf("set project policy-check: %w", err)
	}

	return nil
}
