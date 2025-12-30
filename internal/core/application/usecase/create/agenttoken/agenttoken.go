package agenttoken

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	InitializeLikeUserWithRetry(ctx context.Context, username, password string) error
	CreateAgentToken(ctx context.Context, agentName string) (string, error)
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

	err := u.aiAdapter.InitializeLikeUserWithRetry(ctx, login, password)
	if err != nil {
		return fmt.Errorf("initialize like user: %w", err)
	}

	token, err := u.aiAdapter.CreateAgentToken(ctx, agentName)
	if err != nil {
		return fmt.Errorf("create agent token: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "agent token '%v' created", agentName)
	u.cliAdapter.ReturnText(ctx, token)

	return nil
}
