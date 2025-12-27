package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportPlain struct {
	*cobra.Command
}

type UseCaseGetScanReportPlain interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportPlainCmd(uc UseCaseGetScanReportPlain) CmdGetScanReportPlain {
	cmd := &cobra.Command{
		Use:   "plain <scan-id>",
		Short: "Get scan report plain",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.PlainReport, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report plain' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportPlain{cmd}
}
