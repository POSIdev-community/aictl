package delete

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "Delete resources",
		PersistentPreRunE: _utils.UpdateConfig(cfg),
	}

	cmd.AddCommand(NewDeleteProjectsCommand(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
