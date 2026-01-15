package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanSbom struct {
	*cobra.Command
}

type UseCaseGetScanSbom interface {
	Execute(ctx context.Context, scanId uuid.UUID, outPath string) error
}

func NewGetScanSbomCmd(uc UseCaseGetScanSbom) CmdGetScanSbom {
	var (
		outPath             string
		forceRewriteOutPath bool
	)

	cmd := &cobra.Command{
		Use:   "sbom <scan-id>",
		Short: "Get scan sbom",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outPath != "" {
				if fshelper.PathExists(outPath) && !forceRewriteOutPath {
					return errs.NewValidationError("'output' path exists")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, outPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan sbom' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Destination path for the report file")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Force rewrite output file")

	return CmdGetScanSbom{cmd}
}
