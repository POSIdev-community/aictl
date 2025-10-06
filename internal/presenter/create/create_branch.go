package create

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewCreateBranchCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		projectIdFlag string
	)

	var (
		projectId  uuid.UUID
		branchName string
		scanTarget string
	)

	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Create branch",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			projectId, err = uuid.Parse(projectIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(projectIdFlag, "invalid uuid")
			}

			if scanTarget != "" {
				if !fshelper.PathExists(scanTarget) {
					return errs.NewValidationError(fmt.Sprintf("scan-target path '%s' not exists", scanTarget))
				}
			}

			branchName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.CreateBranchUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, projectId, branchName, scanTarget); err != nil {
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
