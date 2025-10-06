package get

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewGetScanStateCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Get scan stage",
		Use:   "stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScanStateUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get scan state useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	return cmd
}
