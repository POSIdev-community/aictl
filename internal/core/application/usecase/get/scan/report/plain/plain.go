package plain

import (
	"context"
	"fmt"
	"io"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/application/usecase/.utils"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetTemplateId(ctx context.Context, reportType string) (uuid.UUID, error)
	GetReport(ctx context.Context, projectId, scanResultId, templateId uuid.UUID, includeComments, includeDFD, includeGlossary bool, l10n string) (io.ReadCloser, error)
}

type CLI interface {
	ShowReader(r io.Reader) error
	ShowTextf(ctx context.Context, format string, args ...any)
	ShowText(ctx context.Context, text string)
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

	return &UseCase{
		aiAdapter:  aiAdapter,
		cliAdapter: cliAdapter,
		cfg:        cfg,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	u.cliAdapter.ShowTextf(ctx, "getting plain scan report, scan-id '%v'", scanId.String())

	templateId, err := u.aiAdapter.GetTemplateId(ctx, report.PlainReportType)
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

	u.cliAdapter.ShowText(ctx, "plain scan report got")

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
