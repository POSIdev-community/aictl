package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CmdGetScanLogs struct {
	*cobra.Command
}

type UseCaseGetScanLogs interface {
	Execute(ctx context.Context, scanId uuid.UUID, outputPath string) error
}

func NewGetScanLogsCmd(uc UseCaseGetScanLogs) CmdGetScanLogs {
	var (
		outPath             string
		forceRewriteOutPath bool
	)

	cmd := &cobra.Command{
		Use:   "logs <scan-id>",
		Short: "Get scan logs",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outPath == "" {
				return errs.NewValidationRequiredError("output")
			}

			if fshelper.PathExists(outPath) && !forceRewriteOutPath {
				return errs.NewValidationError("'output' path exists")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, outPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan logs' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Force rewrite output file")

	return CmdGetScanLogs{cmd}
}
