package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanAiproj struct {
	*cobra.Command
}

type UseCaseGetScanAiproj interface {
	Execute(ctx context.Context, scanId uuid.UUID, outputPath string) error
}

func NewGetScanAiprojCmd(uc UseCaseGetScanAiproj) CmdGetScanAiproj {
	var outPath string

	cmd := &cobra.Command{
		Short: "Get scan aiproj",
		Use:   "aiproj",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outPath != "" {
				if fshelper.PathExists(outPath) {
					return errs.NewValidationError("'output' path exists")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, outPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan airpoj' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output path")

	return CmdGetScanAiproj{cmd}
}
