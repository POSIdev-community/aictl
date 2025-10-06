package create

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewCreateProjectCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var projectName string

	cmd := &cobra.Command{
		Use:   "project",
		Short: "Create project",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			projectName = args[0]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.CreateProjectUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("create project useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, projectName); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("create project useCase execute: %w", err)
			}

			return nil
		},
	}

	return cmd
}
