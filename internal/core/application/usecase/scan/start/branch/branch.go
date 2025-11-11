package branch

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	StartScanBranch(ctx context.Context, branchId uuid.UUID) (uuid.UUID, error)
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

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config) error {
	scanResultId, err := u.aiAdapter.StartScanBranch(ctx, cfg.BranchId())
	if err != nil {
		return err
	}

	u.cliAdapter.ReturnText(scanResultId.String())

	return nil
}
