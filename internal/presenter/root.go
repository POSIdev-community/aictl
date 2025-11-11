package presenter

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/internal/presenter/context"
	"github.com/POSIdev-community/aictl/internal/presenter/create"
	"github.com/POSIdev-community/aictl/internal/presenter/delete"
	"github.com/POSIdev-community/aictl/internal/presenter/get"
	"github.com/POSIdev-community/aictl/internal/presenter/scan"
	"github.com/POSIdev-community/aictl/internal/presenter/set"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/POSIdev-community/aictl/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var versionFlag bool

	rootCmd := &cobra.Command{
		Use:               "aictl",
		Short:             "Application Inspector ConTroL tool",
		Long:              `aictl - api клиент для Application Inspector`,
		PersistentPreRunE: _utils.InitializeLogger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				l := logger.FromContext(cmd.Context())
				l.StdOut(version.GetVersion())

				return nil
			}

			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					return fmt.Errorf("get help: %w", err)
				}
				return nil
			}

			return nil
		},
	}

	rootCmd.AddCommand(context.NewContextCmd(cfg, depsContainer))
	rootCmd.AddCommand(create.NewCreateCmd(cfg, depsContainer))
	rootCmd.AddCommand(delete.NewDeleteCmd(cfg, depsContainer))
	rootCmd.AddCommand(get.NewGetCmd(cfg, depsContainer))
	rootCmd.AddCommand(scan.NewScanCmd(cfg, depsContainer))
	rootCmd.AddCommand(set.NewSetCmd(cfg, depsContainer))
	rootCmd.AddCommand(update.NewUpdateCmd(cfg, depsContainer))

	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "show version")

	return rootCmd
}
