package aiproj

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	domainaiproj "github.com/POSIdev-community/aictl/internal/core/domain/aiproj"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type CLI interface {
	ShowTextf(ctx context.Context, format string, a ...any)
	ReturnText(ctx context.Context, text string)
}

type UseCase struct {
	cliAdapter CLI
}

func NewUseCase(cliAdapter CLI) (*UseCase, error) {
	if cliAdapter == nil {
		return nil, validation.NewRequiredError("cliAdapter")
	}

	return &UseCase{cliAdapter: cliAdapter}, nil
}

func (u *UseCase) Execute(ctx context.Context, raw []byte, schemaVersion string, jsonOut bool) error {
	report, err := domainaiproj.Check(raw, schemaVersion)
	if err != nil {
		return fmt.Errorf("check aiproj: %w", err)
	}

	if report.OK {
		if jsonOut {
			payload, err := marshalReport(report)
			if err != nil {
				return err
			}

			u.cliAdapter.ReturnText(ctx, payload)

			return nil
		}

		u.cliAdapter.ShowTextf(ctx, "aiproj ok (version %s)", report.Version)

		return nil
	}

	if jsonOut {
		payload, err := marshalReport(report)
		if err != nil {
			return err
		}

		u.cliAdapter.ReturnText(ctx, payload)

		return apperror.NewFailError(report.Kind.ShortMessage())
	}

	return apperror.NewFailError(report.HumanMessage())
}

func marshalReport(report domainaiproj.Report) (string, error) {
	if report.Errors == nil {
		report.Errors = []domainaiproj.Issue{}
	}

	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal check report: %w", err)
	}

	return string(b), nil
}
