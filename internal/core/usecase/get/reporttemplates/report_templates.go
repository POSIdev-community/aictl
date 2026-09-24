package reporttemplates

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type AI interface {
	InitializeWithRetry(ctx context.Context) error
	GetReportTemplates(ctx context.Context, localization string) ([]report.Template, error)
}

type CLI interface {
	ShowReportTemplates(ctx context.Context, templates []report.Template)
	ShowReportTemplatesQuiet(ctx context.Context, templates []report.Template)
}

type UseCase struct {
	aiAdapter  AI
	cliAdapter CLI
}

func NewUseCase(aiAdapter AI, cliAdapter CLI) (*UseCase, error) {
	if aiAdapter == nil {
		return nil, validation.NewRequiredError("aiAdapter")
	}

	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{aiAdapter, cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, filter regexfilter.RegexFilter, quiet bool, localization string) error {
	err := u.aiAdapter.InitializeWithRetry(ctx)
	if err != nil {
		return fmt.Errorf("initialize with retry: %w", err)
	}

	templates, err := u.aiAdapter.GetReportTemplates(ctx, localization)
	if err != nil {
		return fmt.Errorf("get report templates: %w", err)
	}

	filtered := make([]report.Template, 0, len(templates))
	if filter.Empty() {
		filtered = templates
	} else {
		for _, t := range templates {
			if filter.Execute(t.Name) {
				filtered = append(filtered, t)
			}
		}
	}

	if quiet {
		u.cliAdapter.ShowReportTemplatesQuiet(ctx, filtered)
	} else {
		u.cliAdapter.ShowReportTemplates(ctx, filtered)
	}

	return nil
}
