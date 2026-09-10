package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportSarif struct {
	*cobra.Command
}

type UseCaseGetScanReportSarif interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func NewGetScanReportSarifCmd(uc UseCaseGetScanReportSarif) CmdGetScanReportSarif {
	cmd := &cobra.Command{
		Use:   "sarif <scan-id>",
		Short: "Get scan report in SARIF format",
		Long:  `Download the scan report in SARIF format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan report sarif <scan-id> -o ./out.sarif
  aictl get scan report sarif <scan-id> -o ./out.sarif -f`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Sarif, outPath, includeComments, includeDFD, includeGlossary, l10n, report.EmptyFilters()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report sarif' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportSarif{cmd}
}
