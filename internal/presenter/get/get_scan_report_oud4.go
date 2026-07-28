package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportOud4 struct {
	*cobra.Command
}

type UseCaseGetScanReportOud4 interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportOud4Cmd(uc UseCaseGetScanReportOud4) CmdGetScanReportOud4 {
	cmd := &cobra.Command{
		Use:   "oud4 <scan-id>",
		Short: "Get scan report oud4",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Oud4, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report oud4' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportOud4{cmd}
}
