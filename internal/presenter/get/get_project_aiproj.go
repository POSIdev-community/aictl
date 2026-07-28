package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
)

type CmdGetProjectAiproj struct {
	*cobra.Command
}

type UseCaseGetProjectAiproj interface {
	Execute(ctx context.Context, outputPath string) error
}

func NewGetProjectAiprojCmd(uc UseCaseGetProjectAiproj) CmdGetProjectAiproj {
	var (
		outPath             string
		forceRewriteOutPath bool
	)

	cmd := &cobra.Command{
		Use:   "aiproj",
		Short: "Get project aiproj",
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

			if err := uc.Execute(ctx, outPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get project aiproj' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "", "Output path")
	cmd.Flags().BoolVarP(&forceRewriteOutPath, "force", "f", false, "Force rewrite output file")

	return CmdGetProjectAiproj{cmd}
}
