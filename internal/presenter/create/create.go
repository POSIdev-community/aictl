package create

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewCreateCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Create resource",
		PersistentPreRunE: _utils.UpdateConfig(cfg),
	}

	cmd.AddCommand(NewCreateProjectCommand(cfg, depsContainer))
	cmd.AddCommand(NewCreateBranchCommand(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
