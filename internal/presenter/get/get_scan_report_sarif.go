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
	Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string, includeComments, includeDFD, includeGlossary bool) error
}

func NewGetScanReportSarifCmd(uc UseCaseGetScanReportSarif) CmdGetScanReportSarif {
	cmd := &cobra.Command{
		Short: "Get scan report sarif",
		Use:   "sarif",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, destPath, includeComments, includeDFD, includeGlossary); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot sarif: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportSarif{cmd}
}
