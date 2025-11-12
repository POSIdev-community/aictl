package sources

import (
	"context"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	UpdateSources(ctx context.Context, projectId, branchId uuid.UUID, sourcePath string) error
}

type CLI interface {
	ShowText(text string)
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

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, sourcePath string) error {
	u.cliAdapter.ShowText("updating sources")

	err := u.aiAdapter.UpdateSources(ctx, cfg.ProjectId(), cfg.BranchId(), sourcePath)
	if err != nil {
		return err
	}

	u.cliAdapter.ShowText("sources updated")

	return nil
}
