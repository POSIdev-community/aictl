package scan

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdScanStartBranch struct {
	*cobra.Command
}

type UseCaseScanStartBranch interface {
	Execute(ctx context.Context) error
}

func NewScanStartBranchCmd(cfg *config.Config, uc UseCaseScanStartBranch) CmdScanStartBranch {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Start branch scan",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
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

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("scan start: %w", err)
			}

			return nil
		},
	}

	return CmdScanStartBranch{cmd}
}
