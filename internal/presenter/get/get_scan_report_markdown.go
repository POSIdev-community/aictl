package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportMarkdown struct {
	*cobra.Command
}

type UseCaseGetScanReportMarkdown interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportMarkdownCmd(uc UseCaseGetScanReportMarkdown) CmdGetScanReportMarkdown {
	cmd := &cobra.Command{
		Use:   "markdown <scan-id>",
		Short: "Get scan report markdown",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Markdown, destPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report markdown' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportMarkdown{cmd}
}
