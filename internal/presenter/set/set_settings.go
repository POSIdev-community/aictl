package set

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewSetSettingsCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.SetSettingsUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter set settings useCase error: %w", err)
			}

			if err := useCase.Execute(ctx); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter set settings: %w", err)
			}

			return nil
		},
		Short: "Set settings",
		Use:   "settings",
	}

	return cmd
}
