package reports

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/google/uuid"

	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type CLI interface {
}

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetTemplateId(ctx context.Context, reportType string) (uuid.UUID, error)
	GetProject(ctx context.Context, id uuid.UUID) (*project.Project, error)
	GetBranches(ctx context.Context, projectId uuid.UUID) (branches []branch.Branch, err error)
	GetLastScan(ctx context.Context, branchId uuid.UUID) (*scan.Scan, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool) (io.ReadCloser, error)
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

func (u *UseCase) Execute(ctx context.Context, projectIds []uuid.UUID, sarif bool, plain bool, destPath string, includeComments, includeDFD, includeGlossary bool) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("could not initialize with jwt retry: %w", err)
	}

	var reportType string
	switch {
	case sarif:
		reportType = report.SarifReportType
	case plain:
		reportType = report.PlainReportType
	}
	templateId, err := u.aiAdapter.GetTemplateId(ctx, reportType)
	if err != nil {
		return fmt.Errorf("get reports usecase: %w", err)
	}

	for _, projectId := range projectIds {
		proj, err := u.aiAdapter.GetProject(ctx, projectId)
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

		file, err := u.aiAdapter.GetReport(ctx, projectId, scanResult.Id, templateId, includeComments, includeDFD, includeGlossary)
		if err != nil {
			return fmt.Errorf("get reports usecase get report: %w", err)
		}

		err = utils.CopyFileToPath(file, filepath.Join(destPath, proj.Name+".sarif"))
		if err != nil {
			return fmt.Errorf("get reports usecase: %w", err)
		}
	}

	return nil
}
