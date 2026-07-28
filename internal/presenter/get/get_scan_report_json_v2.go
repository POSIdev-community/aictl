package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

type CmdGetScanReportJsonV2 struct {
	*cobra.Command
}

type UseCaseGetScanReportJsonV2 interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportJsonV2Cmd(uc UseCaseGetScanReportJsonV2) CmdGetScanReportJsonV2 {
	cmd := &cobra.Command{
		Use:   "json-v2 <scan-id>",
		Short: "Get scan report json v2",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.JsonV2, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report json-v2' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportJsonV2{cmd}
}
