package set

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func NewSetProjectSettingsCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.SetProjectSettingsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter set project settings useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg.ProjectId(), uuid.New()); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter set project settings: %w", err)
			}

			return nil
		},
		Short: "Set project settings",
		Use:   "settings",
	}

	return cmd
}
