package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportSans struct {
	*cobra.Command
}

type UseCaseGetScanReportSans interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportSansCmd(uc UseCaseGetScanReportSans) CmdGetScanReportSans {
	cmd := &cobra.Command{
		Use:   "sans <scan-id>",
		Short: "Get scan report sans",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Sans, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report sans' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportSans{cmd}
}
