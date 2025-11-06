package get

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/pkg/errs"
)

var (
	projectIdFlag string
	scanId        uuid.UUID
)

func NewGetScanCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Get scan",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			if len(args) < 1 {
				return errs.NewValidationError("missing scan id")
			}

			scanIdFlag := args[0]

			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return errs.NewValidationFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase, err := depsContainer.GetScanUseCase(ctx, cfg)
			if err != nil {
				return fmt.Errorf("get projects useCase error: %w", err)
			}

			if err := useCase.Execute(ctx, cfg, scanId); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("get projects: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(NewGetScanAiprojCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanLogsCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanReportCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanSbomCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanStateCmd(cfg, depsContainer))

	cmd.PersistentFlags().StringVarP(&projectIdFlag, "project-id", "p", "", "project id")

	return cmd
}
