package scan

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewScanCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "scan",
		Short:             "Scan ",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(NewScanAwaitCommand(cfg, depsContainer))
	cmd.AddCommand(NewScanStartCommand(cfg, depsContainer))
	cmd.AddCommand(NewScanStopCommand(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
