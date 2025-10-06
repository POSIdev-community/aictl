package update

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"strings"
)

func NewUpdateSourcesCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		path      string
		projectId uuid.UUID
		branchId  uuid.UUID
	)

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "Update sources",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			path = strings.TrimSpace(args[0])
			if path == "" {
				return errs.NewValidationError("empty sources path")
			}

			if !fshelper.PathExists(path) {
				return errs.NewValidationError("path does not exist")
			}

			var err error

			projectId, err = uuid.Parse(projectIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(projectIdFlag, "invalid uuid")
			}

			branchId, err = uuid.Parse(branchIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(projectIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.UpdateSourcesUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("update sources useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, projectId, branchId, path); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "branch id")

	cmd.AddCommand(NewUpdateSourcesGitCommand(cfg, depsContainer))

	return cmd
}
