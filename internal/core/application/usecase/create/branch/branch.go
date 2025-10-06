package branch

import (
	"context"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/port"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
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

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, projectId uuid.UUID, branchName, scanTarget string) error {
	branchId, err := u.aiAdapter.CreateBranch(ctx, projectId, branchName, scanTarget)
	if err != nil {
		return fmt.Errorf("usecase create branch: %w", err)
	}

	u.cliAdapter.ShowText(branchId.String())

	return nil
}
