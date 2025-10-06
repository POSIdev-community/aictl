package get

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewGetScanReportSarifCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Get scan report sarif",
		Use:   "sarif",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScanReportSarifUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter get scan repot sarif useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId, destPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot sarif: %w", err)
			}

			return nil
		},
	}

	return cmd
}
