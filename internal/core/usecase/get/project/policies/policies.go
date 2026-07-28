package policies

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetProjectPolicies(ctx context.Context, projectId uuid.UUID) (io.ReadCloser, error)
}

type CLI interface {
	ShowReader(r io.Reader) error
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	r, err := u.aiAdapter.GetProjectPolicies(ctx, u.cfg.ProjectId())
	if err != nil {
		return fmt.Errorf("get project policies: %w", err)
	}
	defer func() {
		_ = r.Close()
	}()

	if err := u.cliAdapter.ShowReader(r); err != nil {
		return fmt.Errorf("print project policies: %w", err)
	}

	return nil
}
