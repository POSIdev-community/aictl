package get

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
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
		Short: "Get scan SBOM",
		Long:  `Download the scan SBOM (Software Bill of Materials). Scan id comes from argument or stdin. Output path via -o; use -f to overwrite.`,
		Example: `  aictl get scan sbom <scan-id> -o ./sbom.json
  aictl get scan sbom <scan-id> -o ./sbom.json -f`,
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

				return fmt.Errorf("'get scan sbom' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output file path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Overwrite existing output file")

	return CmdGetScanSbom{cmd}
}
