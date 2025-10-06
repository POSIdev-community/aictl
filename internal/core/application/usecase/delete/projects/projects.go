package projects

import (
	"context"
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

func (u *UseCase) Execute(ctx context.Context, projectIds []uuid.UUID) error {
	for _, projectId := range projectIds {
		err := u.aiAdapter.DeleteProject(ctx, projectId)
		if err != nil {
			return err
		}

		u.cliAdapter.ShowText(projectId.String())
	}

	return nil
}
