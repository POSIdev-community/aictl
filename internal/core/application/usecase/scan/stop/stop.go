package stop

import (
	"context"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	StopScan(ctx context.Context, scanResultId uuid.UUID) error
}

type CLI interface {
	ReturnText(string)
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

func (u *UseCase) Execute(ctx context.Context, scanResultId uuid.UUID) error {
	err := u.aiAdapter.StopScan(ctx, scanResultId)
	if err != nil {
		return err
	}

	u.cliAdapter.ReturnText(scanResultId.String())

	return nil
}
