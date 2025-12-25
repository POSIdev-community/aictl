package branch

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	StartScanBranch(ctx context.Context, branchId uuid.UUID, scanLabel string) (uuid.UUID, error)
}

type CLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
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

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanLabel string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "starting scan, project-id '%v', branch-id '%v'", u.cfg.ProjectId(), u.cfg.BranchId())

	scanResultId, err := u.aiAdapter.StartScanBranch(ctx, u.cfg.BranchId(), scanLabel)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowTextf(ctx, "scan started, scan-id '%v'", scanResultId)
	u.cliAdapter.ReturnText(ctx, scanResultId.String())

	return nil
}
