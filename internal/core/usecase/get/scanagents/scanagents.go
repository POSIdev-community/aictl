package scanagents

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error)
}

type CLI interface {
	ShowScanAgents(ctx context.Context, agents []scanagent.ScanAgent)
	ShowScanAgentsQuiet(ctx context.Context, agents []scanagent.ScanAgent)
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

func (u *UseCase) Execute(ctx context.Context, quiet bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	agents, err := u.aiAdapter.GetScanAgents(ctx)
	if err != nil {
		return fmt.Errorf("get scan agents: %w", err)
	}

	if quiet {
		u.cliAdapter.ShowScanAgentsQuiet(ctx, agents)
	} else {
		u.cliAdapter.ShowScanAgents(ctx, agents)
	}

	return nil
}
