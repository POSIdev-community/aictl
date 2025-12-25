package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportSarif struct {
	*cobra.Command
}

type UseCaseGetScanReportSarif interface {
	Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string, includeComments, includeDFD, includeGlossary bool, l10n string) error
}

func NewGetScanReportSarifCmd(uc UseCaseGetScanReportSarif) CmdGetScanReportSarif {
	cmd := &cobra.Command{
		Use:   "sarif <scan-id>",
		Short: "Get scan report sarif",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, destPath, includeComments, includeDFD, includeGlossary, l10n); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan report sarif' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportSarif{cmd}
}
