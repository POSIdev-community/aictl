package create

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/spf13/cobra"
)

type CmdCreateBranch struct {
	*cobra.Command
}

type UseCaseCreateBranch interface {
	Execute(ctx context.Context, branchName, scanTarget string, safe bool) error
}

func NewCreateBranchCmd(cfg *config.Config, uc UseCaseCreateBranch) CmdCreateBranch {

	var (
		projectIdFlag string
	)

	var (
		branchName string
		scanTarget string
	)

	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Create branch",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			if err := cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			if scanTarget != "" {
				if !fshelper.PathExists(scanTarget) {
					return errs.NewValidationError(fmt.Sprintf("scan-target path '%s' not exists", scanTarget))
				}
			}

			args = _utils.ReadArgsFromStdin(args)
			branchName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, branchName, scanTarget, safeFlag); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&scanTarget, "scan-target", "s", "", "scan target path")

	return CmdCreateBranch{cmd}
}
