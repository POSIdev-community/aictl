package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportJson struct {
	*cobra.Command
}

type UseCaseGetScanReportJson interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string, filters report.Filters) error
}

func NewGetScanReportJsonCmd(uc UseCaseGetScanReportJson) CmdGetScanReportJson {
	cmd := &cobra.Command{
		Use:   "json <scan-id>",
		Short: "Get scan report in JSON format",
		Long:  `Download the scan report in JSON format for the given scan id. Project id comes from context or parent -p. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan report json <scan-id> -o ./out.json
  aictl get scan report json <scan-id> -o ./out.json -f`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Json, outPath, includeComments, includeDFD, includeGlossary, l10n, report.EmptyFilters()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report json' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportJson{cmd}
}
