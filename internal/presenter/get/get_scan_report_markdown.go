package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
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
		Short: "Get scan report in Markdown format",
		Long:  `Download the scan report in Markdown format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan report markdown <scan-id> -o ./out.md
  aictl get scan report markdown <scan-id> -o ./out.md -f`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Markdown, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report markdown' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportMarkdown{cmd}
}
