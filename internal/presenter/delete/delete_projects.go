package delete

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewDeleteProjectsCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var projectIds []uuid.UUID

	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Delete AI projects",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			args = _utils.ReadArgsFromStdin(args)
			projectIds, err = _utils.ParseUUIDs(args)
			if err != nil {
				return fmt.Errorf("get reports project ids parse error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.DeleteProjectsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, projectIds); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	return cmd
}
