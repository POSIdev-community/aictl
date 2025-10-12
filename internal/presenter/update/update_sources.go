package update

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/fshelper"
	"github.com/spf13/cobra"
	"strings"
)

func NewUpdateSourcesCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		path string
	)

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "Update sources",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			var err error

			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			if err = cfg.UpdateBranchId(branchIdFlag); err != nil {
				return err
			}

			path = strings.TrimSpace(args[0])
			if path == "" {
				return errs.NewValidationError("empty sources path")
			}

			if !fshelper.PathExists(path) {
				return errs.NewValidationError("path does not exist")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.UpdateSourcesUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("update sources useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, path); err != nil {
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
