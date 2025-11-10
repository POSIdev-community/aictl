package state

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	GetScanStage(ctx context.Context, projectId, scanId uuid.UUID) (scanstage.ScanStage, error)
}

type CLI interface {
	ReturnText(text string)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
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

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, scanId uuid.UUID) error {
	scanStage, err := u.aiAdapter.GetScanStage(ctx, cfg.ProjectId(), scanId)
	if err != nil {
		return err
	}

	u.cliAdapter.ReturnText(scanStage.Stage)

	return nil
}
