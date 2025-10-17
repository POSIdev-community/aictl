package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

func NewSetCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "set",
		Short:             "Set",
		PersistentPreRunE: _utils.UpdateConfig(cfg),
	}

	cmd.AddCommand(NewSetProjectCmd(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
