package logs

import (
	"context"
	"fmt"
	"io"

	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScanLogs(ctx context.Context, projectId, scanResultId uuid.UUID) (io.ReadCloser, error)
}

type CLI interface {
	ShowTextf(ctx context.Context, format string, args ...any)
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

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "getting scan logs, scan-id '%v'", scanId.String())

	r, err := u.aiAdapter.GetScanLogs(ctx, u.cfg.ProjectId(), scanId)
	if err != nil {
		return fmt.Errorf("get scan logs: %w", err)
	}

	defer func() {
		_ = r.Close()
	}()

	u.cliAdapter.ShowTextf(ctx, "scan logs got")

	if err := utils.CopyFileToPath(r, fullDestPath); err != nil {
		return fmt.Errorf("copy logs to path %s: %w", fullDestPath, err)
	}

	return nil
}
