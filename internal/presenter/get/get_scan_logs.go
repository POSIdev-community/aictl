package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
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
		Long:  `Download scan logs to a file. Scan id comes from argument or stdin. Output path -o is required; use -f to overwrite.`,
		Example: `  aictl get scan logs <scan-id> -o ./scan.log
  aictl get scan logs <scan-id> -o ./scan.log -f`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outPath == "" {
				return validation.NewRequiredError("output")
			}

			if fshelper.PathExists(outPath) && !forceRewriteOutPath {
				return validation.NewError("'output' path exists")
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

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output file path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Overwrite existing output file")

	return CmdGetScanLogs{cmd}
}
