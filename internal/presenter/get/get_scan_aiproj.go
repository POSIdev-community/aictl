package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdGetScanAiproj struct {
	*cobra.Command
}

type UseCaseGetScanAiproj interface {
	Execute(ctx context.Context, scanId uuid.UUID, outputPath string) error
}

func NewGetScanAiprojCmd(uc UseCaseGetScanAiproj) CmdGetScanAiproj {
	var (
		outPath             string
		forceRewriteOutPath bool
	)

	cmd := &cobra.Command{
		Use:   "aiproj <scan-id>",
		Short: "Get scan aiproj",
		Long:  `Download the .aiproj configuration used for the scan. Scan id comes from argument or stdin. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan aiproj <scan-id> -o ./scan.aiproj
  aictl get scan aiproj <scan-id> -o ./scan.aiproj -f`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outPath != "" {
				if fshelper.PathExists(outPath) && !forceRewriteOutPath {
					return validation.NewError("'output' path exists")
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

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output file path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Overwrite existing output file")

	return CmdGetScanAiproj{cmd}
}
