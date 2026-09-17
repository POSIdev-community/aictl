package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportPlain struct {
	*cobra.Command
}

type UseCaseGetScanReportPlain interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func NewGetScanReportPlainCmd(uc UseCaseGetScanReportPlain) CmdGetScanReportPlain {
	cmd := &cobra.Command{
		Use:   "plain <scan-id>",
		Short: "Get scan report in plain text format",
		Long:  `Download the scan report in plain text format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan report plain <scan-id> -o ./out.txt
  aictl get scan report plain <scan-id> -o ./out.txt -f`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.PlainReport, outPath, includeComments, includeDFD, includeGlossary, l10n, report.EmptyFilters()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report plain' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportPlain{cmd}
}
