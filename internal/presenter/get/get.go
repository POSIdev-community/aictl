package get

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewGetCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "get",
		Short:             "Get resources",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(NewGetProjectsCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetReportsCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScanCmd(cfg, depsContainer))
	cmd.AddCommand(NewGetScansCmd(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
