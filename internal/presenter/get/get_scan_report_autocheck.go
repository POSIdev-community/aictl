package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportAutocheck struct {
	*cobra.Command
}

type UseCaseGetScanReportAutocheck interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportAutocheckCmd(uc UseCaseGetScanReportAutocheck) CmdGetScanReportAutocheck {
	cmd := &cobra.Command{
		Use:   "autocheck <scan-id>",
		Short: "Get scan report autocheck",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.AutoCheck, destPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report autocheck' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportAutocheck{cmd}
}
