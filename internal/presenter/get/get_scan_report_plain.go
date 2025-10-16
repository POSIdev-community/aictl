package get

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
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

			reportFilePath := fileName
			if reportFilePath != "" {
				reportFilePath = filepath.Join(destPath, reportFilePath)
			}

			if err := useCase.Execute(ctx, cfg, scanId, reportFilePath); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("presenter get scan repot plain: %w", err)
			}

			return nil
		},
	}

	return cmd
}
