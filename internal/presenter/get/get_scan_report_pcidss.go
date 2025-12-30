package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportPcidss struct {
	*cobra.Command
}

type UseCaseGetScanReportPcidss interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportPcidssCmd(uc UseCaseGetScanReportPcidss) CmdGetScanReportPcidss {
	cmd := &cobra.Command{
		Use:   "pcidss <scan-id>",
		Short: "Get scan report pcidss",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Pcidss, destPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report pcidss' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportPcidss{cmd}
}
