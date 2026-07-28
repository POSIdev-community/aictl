package checkpolicies

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/policystate"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanPolicyState(ctx context.Context, projectId, scanId uuid.UUID) (policystate.State, error)
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, failOnPoliciesRejected bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	state, err := u.aiAdapter.GetScanPolicyState(ctx, u.cfg.ProjectId(), scanId)
	if err != nil {
		return fmt.Errorf("check scan policies: %w", err)
	}

	u.cliAdapter.ReturnText(ctx, state.String())

	if failOnPoliciesRejected && state.IsRejected() {
		return apperror.NewPolicyFailError(state.String())
	}

	return nil
}
