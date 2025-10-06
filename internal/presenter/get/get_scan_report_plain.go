package get

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewGetScanReportPlainCmd(cfg *config.Config, depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Get scan report plain",
		Use:   "plain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScanReportPlainUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("presenter get scan repot plain useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId, destPath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot plain: %w", err)
			}

			return nil
		},
	}

	return cmd
}
