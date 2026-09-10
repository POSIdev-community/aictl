package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportSans struct {
	*cobra.Command
}

type UseCaseGetScanReportSans interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func NewGetScanReportSansCmd(uc UseCaseGetScanReportSans) CmdGetScanReportSans {
	cmd := &cobra.Command{
		Use:   "sans <scan-id>",
		Short: "Get scan report in SANS format",
		Long:  `Download the scan report in SANS format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan report sans <scan-id> -o ./out.xml
  aictl get scan report sans <scan-id> -o ./out.xml -f`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Sans, outPath, includeComments, includeDFD, includeGlossary, l10n, report.EmptyFilters()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report sans' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportSans{cmd}
}
