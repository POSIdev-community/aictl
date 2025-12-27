package healthcheck

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetHealthcheck(ctx context.Context) (bool, error)
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
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
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

	h, err := u.aiAdapter.GetHealthcheck(ctx)
	if err != nil {
		return fmt.Errorf("get healthcheck: %w", err)
	}

	var text string
	if h {
		text = "Healthy"
	} else {
		text = "Unhealthy"
	}

	u.cliAdapter.ReturnText(ctx, text)

	return nil
}
