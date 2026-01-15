package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportGitlab struct {
	*cobra.Command
}

type UseCaseGetScanReportGitlab interface {
	Execute(ctx context.Context, scanId uuid.UUID, reportType report.ReportType, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportGitlabCmd(uc UseCaseGetScanReportGitlab) CmdGetScanReportGitlab {
	cmd := &cobra.Command{
		Use:   "gitlab <scan-id>",
		Short: "Get scan report gitlab",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, report.Gitlab, outPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report gitlab' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportGitlab{cmd}
}
