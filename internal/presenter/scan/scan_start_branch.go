package scan

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewScanStartBranchCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Start branch scan",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			args = _utils.ReadArgsFromStdin(args)
			var branchIdFlag string
			if len(args) > 0 {
				branchIdFlag = args[0]
			}

			var err error
			if branchIdFlag != "" {
				err = cfg.SetProjectId(branchIdFlag)
				if err != nil {
					return err
				}
			} else {
				err = cfg.Validate(false, true)
				if err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.ScanStartBranchUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("scan start: %w", err)
			}

			return nil
		},
	}

	return cmd
}
