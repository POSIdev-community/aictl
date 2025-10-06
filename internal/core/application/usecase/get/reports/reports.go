package reports

import (
	"context"
	"fmt"
	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
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

func (u *UseCase) Execute(ctx context.Context, projectIds []uuid.UUID, sarif bool, plain bool, destPath string) error {
	var reportType string
	switch {
	case sarif:
		reportType = port.SarifReportType
	case plain:
		reportType = port.PlainReportType
	}
	templateId, err := u.aiAdapter.GetTemplateId(ctx, reportType)
	if err != nil {
		return fmt.Errorf("get reports usecase: %w", err)
	}

	for _, projectId := range projectIds {
		project, err := u.aiAdapter.GetProject(ctx, projectId)
		if err != nil {
			return fmt.Errorf("get reports usecase get project: %w", err)
		}

		branches, err := u.aiAdapter.GetBranches(ctx, projectId)
		if err != nil {
			return fmt.Errorf("get reports usecase get project branches: %w", err)
		}

		branchId := branches[0].Id

		scanResult, err := u.aiAdapter.GetLastScan(ctx, branchId)
		if err != nil {
			return fmt.Errorf("get reports usecase get project scanResult: %w", err)
		}

		file, err := u.aiAdapter.GetReport(ctx, projectId, scanResult.Id, templateId)
		if err != nil {
			return fmt.Errorf("get reports usecase get report: %w", err)
		}

		err = utils.CopyFileToPath(file, destPath, project.Name+".sarif")
		if err != nil {
			return fmt.Errorf("get reports usecase: %w", err)
		}
	}

	return nil
}
