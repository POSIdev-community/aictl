package create

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/spf13/cobra"
)

func NewCreateBranchCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

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
			//log := logger.FromContext(ctx)
			//
			//log.StdErrF("create branch")

			useCase, err := depsContainer.CreateBranchUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, branchName, scanTarget, safeFlag); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&scanTarget, "scan-target", "s", "", "scan target path")

	return cmd
}
