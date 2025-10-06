package sarif

import (
	"context"
	"fmt"
	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, scanId uuid.UUID, destPath string) error {
	project, err := u.aiAdapter.GetProject(ctx, cfg.ProjectId())
	if err != nil {
		return err
	}

	templateId, err := u.aiAdapter.GetTemplateId(ctx, port.SarifReportType)
	if err != nil {
		return err
	}

	report, err := u.aiAdapter.GetReport(ctx, cfg.ProjectId(), scanId, templateId)

	err = utils.CopyFileToPath(report, destPath, project.Name+".sarif")
	if err != nil {
		return fmt.Errorf("get scan report sarif usecase: %w", err)
	}

	return nil
}
