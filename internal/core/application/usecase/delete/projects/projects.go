package projects

import (
	"context"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	DeleteProject(context context.Context, projectId uuid.UUID) error
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

func (u *UseCase) Execute(ctx context.Context, projectIds []uuid.UUID) error {
	for _, projectId := range projectIds {
		err := u.aiAdapter.DeleteProject(ctx, projectId)
		if err != nil {
			return err
		}

		u.cliAdapter.ReturnText(projectId.String())
	}

	return nil
}
