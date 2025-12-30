package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportJson struct {
	*cobra.Command
}

type UseCaseGetScanReportJson interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportJsonCmd(uc UseCaseGetScanReportJson) CmdGetScanReportJson {
	cmd := &cobra.Command{
		Use:   "json <scan-id>",
		Short: "Get scan report json",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Json, destPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report json' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportJson{cmd}
}
