package update

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewUpdateSourcesGitCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "git",
		Short: "Update sources git",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.UpdateSourcesGitUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("update sources useCase error: %w", err)
			}

			if err := useCase.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	return cmd
}
