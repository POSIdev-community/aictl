package update

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

var (
	projectIdFlag string
	branchIdFlag  string
)

func NewUpdateCmd(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update",
		Short:             "Update resources",
		PersistentPreRunE: _utils.UpdateConfig(cfg),
	}

	cmd.AddCommand(NewUpdateSourcesCommand(cfg, depsContainer))

	_utils.AddConnectionPersistentFlags(cmd)

	return cmd
}
