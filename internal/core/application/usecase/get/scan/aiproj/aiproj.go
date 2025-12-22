package aiproj

import (
	"context"
	"fmt"
	"os"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetScan(ctx context.Context, projectId uuid.UUID, scanId uuid.UUID) (*scan.Scan, error)
	GetScanAiproj(ctx context.Context, projectId uuid.UUID, settingsId uuid.UUID) (string, error)
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

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, outputPath string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	projectId := u.cfg.ProjectId()

	s, err := u.aiAdapter.GetScan(ctx, projectId, scanId)
	if err != nil {
		return err
	}

	aiproj, err := u.aiAdapter.GetScanAiproj(ctx, projectId, s.SettingsId)
	if err != nil {
		return err
	}

	if outputPath == "" {
		u.cliAdapter.ReturnText(ctx, aiproj)
	} else {
		err := os.WriteFile(outputPath, []byte(aiproj), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
