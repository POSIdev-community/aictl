package scans

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	GetScans(ctx context.Context, branchId uuid.UUID) ([]scan.Scan, error)
}

type CLI interface {
	ShowScans(scans []scan.Scan)
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

func (u *UseCase) Execute(ctx context.Context, branchId uuid.UUID) error {
	scans, err := u.aiAdapter.GetScans(ctx, branchId)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowScans(scans)

	return nil
}
