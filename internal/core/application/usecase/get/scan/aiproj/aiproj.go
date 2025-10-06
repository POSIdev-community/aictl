package aiproj

import (
	"context"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"os"
)

type UseCase struct {
	aiAdapter  port.Ai
	cliAdapter port.Cli
}

func NewUseCase(aiAdapter port.Ai, cliAdapter port.Cli) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, scanId uuid.UUID, outputPath string) error {
	projectId := cfg.ProjectId()

	scan, err := u.aiAdapter.GetScan(ctx, projectId, scanId)
	if err != nil {
		return err
	}

	aiproj, err := u.aiAdapter.GetScanAiproj(ctx, projectId, scan.SettingsId)
	if err != nil {
		return err
	}

	if outputPath == "" {
		u.cliAdapter.ShowText(aiproj)
	} else {
		err := os.WriteFile(outputPath, []byte(aiproj), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
