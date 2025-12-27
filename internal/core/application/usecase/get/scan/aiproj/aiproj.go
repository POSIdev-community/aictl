package aiproj

import (
	"context"
	"fmt"
	"io"

	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScan(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (*scan.Scan, error)
	GetScanAiproj(ctx context.Context, projectId, scanSettingsId uuid.UUID) (io.ReadCloser, error)
}

type CLI interface {
	ShowReader(r io.Reader) error
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	s, err := u.aiAdapter.GetScan(ctx, u.cfg.ProjectId(), scanId)
	if err != nil {
		return err
	}

	r, err := u.aiAdapter.GetScanAiproj(ctx, u.cfg.ProjectId(), s.SettingsId)
	if err != nil {
		return err
	}

	defer func() {
		_ = r.Close()
	}()

	u.cliAdapter.ShowTextf(ctx, "aiproj got")

	if fullDestPath != "" {
		if err := utils.CopyFileToPath(r, fullDestPath); err != nil {
			return fmt.Errorf("copy aiproj to path %s: %w", fullDestPath, err)
		}

		return nil
	}

	if err := u.cliAdapter.ShowReader(r); err != nil {
		return fmt.Errorf("print aiproj: %w", err)
	}

	return nil
}
