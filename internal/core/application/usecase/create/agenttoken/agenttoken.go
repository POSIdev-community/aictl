package agenttoken

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	CreateAgentToken(ctx context.Context, login, password, agentName string, tlsSkip bool) (string, error)
}

type CLI interface {
	ReturnText(ctx context.Context, text string)
	ShowTextf(ctx context.Context, format string, a ...any)
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

	if cfg == nil {
		return nil, errs.NewValidationRequiredError("cfg")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, login, password, agentName string) error {
	u.cliAdapter.ShowTextf(ctx, "creating agent token '%v'", agentName)

	token, err := u.aiAdapter.CreateAgentToken(ctx, login, password, agentName, u.cfg.TLSSkip())
	if err != nil {
		return fmt.Errorf("create agent token: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "agent token '%v' created", agentName)
	u.cliAdapter.ReturnText(ctx, token)

	return nil
}
