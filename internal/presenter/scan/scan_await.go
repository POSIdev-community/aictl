package scan

import (
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func NewScanAwaitCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		projectIdFlag string
		scanIdFlag    string
		scanId        uuid.UUID
	)

	cmd := &cobra.Command{
		Use:   "await",
		Short: "Await scan",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			scanIdFlag = args[0]

			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.ScanAwaitUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("scan start: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return cmd
}
