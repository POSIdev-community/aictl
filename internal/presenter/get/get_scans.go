package get

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdGetScans struct {
	*cobra.Command
}

type UseCaseGetScans interface {
	Execute(ctx context.Context) error
}

func NewGetScansCmd(cfg *config.Config, uc UseCaseGetScans) CmdGetScans {
	cmd := &cobra.Command{
		Short: "Get scans",
		Use:   "scans",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			branchIdFlag := args[0]

			if err := cfg.UpdateBranchId(branchIdFlag); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scans' usecase call: %w", err)
			}

			return nil
		},
	}

	return CmdGetScans{cmd}
}
