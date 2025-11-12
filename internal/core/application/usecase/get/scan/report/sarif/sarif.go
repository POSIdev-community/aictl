package sarif

import (
	"context"
	"fmt"
	"io"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"

	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	GetTemplateId(ctx context.Context, reportType string) (uuid.UUID, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool) (io.ReadCloser, error)
}

type CLI interface {
	ShowReader(r io.Reader) error
	ShowTextf(format string, args ...any)
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, cfg *config.Config, scanId uuid.UUID, fullDestPath string, includeComments, includeDFD, includeGlossary bool) error {
	u.cliAdapter.ShowTextf("getting sarif scan report, id '%v'", scanId.String())

	templateId, err := u.aiAdapter.GetTemplateId(ctx, report.SarifReportType)
	if err != nil {
		return err
	}

	report, err := u.aiAdapter.GetReport(ctx, cfg.ProjectId(), scanId, templateId, includeComments, includeDFD, includeGlossary)
	if err != nil {
		return fmt.Errorf("get scan report: %w", err)
	}

	defer func() {
		_ = report.Close()
	}()

	if fullDestPath != "" {
		if err := utils.CopyFileToPath(report, fullDestPath); err != nil {
			return fmt.Errorf("copy report to path %s: %w", fullDestPath, err)
		}

		return nil
	}

	u.cliAdapter.ShowText("sarif scan report got")

	if err := u.cliAdapter.ShowReader(report); err != nil {
		return fmt.Errorf("print report: %w", err)
	}

	return nil
}
