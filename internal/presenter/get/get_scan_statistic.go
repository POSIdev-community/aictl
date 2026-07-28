package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdGetScanStatistic struct {
	*cobra.Command
}

type UseCaseGetScanStatistic interface {
	Execute(ctx context.Context, scanId uuid.UUID, outPath string, json bool) error
}

func NewGetScanStatisticCmd(uc UseCaseGetScanStatistic) CmdGetScanStatistic {
	var (
		outPath             string
		forceRewriteOutPath bool
		json                bool
	)

	cmd := &cobra.Command{
		Use:   "statistic <scan-id>",
		Short: "Get scan statistic",
		Long:  `Download or print scan statistics. Scan id comes from argument or stdin. Output path via -o; use -f to overwrite. Use --json for JSON output.`,
		Example: `  aictl get scan statistic <scan-id>
  aictl get scan statistic <scan-id> --json
  aictl get scan statistic <scan-id> -o ./stat.json -f`,
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

			if err := uc.Execute(ctx, scanId, outPath, json); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scan statistic' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output file path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Overwrite existing output file")
	cmd.Flags().BoolVar(&json, "json", false, "Output in JSON format")

	return CmdGetScanStatistic{cmd}
}
