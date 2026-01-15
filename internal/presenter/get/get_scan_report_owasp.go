package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportOwasp struct {
	*cobra.Command
}

type UseCaseGetScanReportOwasp interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportOwaspCmd(uc UseCaseGetScanReportOwasp) CmdGetScanReportOwasp {
	cmd := &cobra.Command{
		Use:   "owasp <scan-id>",
		Short: "Get scan report owasp",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.AutoCheck, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report owasp' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportOwasp{cmd}
}
