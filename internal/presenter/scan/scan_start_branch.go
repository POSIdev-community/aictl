package scan

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanStartBranch struct {
	*cobra.Command
}

type UseCaseScanStartBranch interface {
	Execute(ctx context.Context, scanLabel string, scanType scantype.Type) error
}

func NewScanStartBranchCmd(cfg *config.Config, uc UseCaseScanStartBranch) CmdScanStartBranch {
	var projectIdFlag string

	cmd := &cobra.Command{
		Use:   "branch <branch-id>",
		Short: "Start branch scan",
		Long:  `Start a scan on a branch. Branch id comes from the argument or context; project id from context or -p.`,
		Example: `  aictl scan start branch <branch-id> -p <project-id>
  aictl scan start branch --scan-label nightly --full-scan`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)

			var branchIdFlag string
			if len(args) > 0 {
				branchIdFlag = args[0]
			}

			if err := cfg.UpdateBranchId(branchIdFlag); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanLabel, scanTypeFromFlags()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan start branch' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")

	return CmdScanStartBranch{cmd}
}
