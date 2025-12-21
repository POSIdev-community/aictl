package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanReportPlain struct {
	*cobra.Command
}

type UseCaseGetScanReportPlain interface {
	Execute(ctx context.Context, scanId uuid.UUID, fullDestPath string, includeComments, includeDFD, includeGlossary bool) error
}

func NewGetScanReportPlainCmd(uc UseCaseGetScanReportPlain) CmdGetScanReportPlain {
	cmd := &cobra.Command{
		Short: "Get scan report plain",
		Use:   "plain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, destPath, includeComments, includeDFD, includeGlossary); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot plain: %w", err)
			}

			return nil
		},
	}

	return CmdGetScanReportPlain{cmd}
}
