package report

import (
	"context"
	"fmt"
	"io"

	utils "github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetCustomTemplateId(ctx context.Context, reportName string) (uuid.UUID, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string) (io.ReadCloser, error)
}

type CLI interface {
	ShowReader(r io.Reader) error
	ShowTextf(ctx context.Context, format string, args ...any)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
	cfg        *config.Config
}

func NewUseCase(aiAdapter AI, cliAdapter CLI, cfg *config.Config) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, errs.NewValidationRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, errs.NewValidationRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter, cfg}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, customReportName string, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "getting '%s' scan report, scan-id '%v'", customReportName, scanId.String())

	templateId, err := u.aiAdapter.GetCustomTemplateId(ctx, customReportName)
	if err != nil {
		return err
	}

	r, err := u.aiAdapter.GetReport(ctx, u.cfg.ProjectId(), scanId, templateId, includeComments, includeDFD, includeGlossary, l10n)
	if err != nil {
		return fmt.Errorf("get scan report: %w", err)
	}

	defer func() {
		_ = r.Close()
	}()

	u.cliAdapter.ShowTextf(ctx, "'%s' scan report got", customReportName)

	if fullDestPath != "" {
		if err := utils.CopyFileToPath(r, fullDestPath); err != nil {
			return fmt.Errorf("copy report to path %s: %w", fullDestPath, err)
		}

		return nil
	}

	if err := u.cliAdapter.ShowReader(r); err != nil {
		return fmt.Errorf("print report: %w", err)
	}

	return nil
}
