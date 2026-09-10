package update

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdUpdate struct {
	*cobra.Command
}

type PersistentPreRunEUpdateCmd _utils.RunE

func NewPersistentPreRunEUpdateCmd(cfg *config.Config) PersistentPreRunEUpdateCmd {
	return _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg))
}

var (
	projectIdFlag string
	branchIdFlag  string
)

func NewUpdateCmd(cfg *config.Config, cmdUpdateSources CmdUpdateSources, cmdUpdateSbom CmdUpdateSbom, cmdUpdateScaFeeds CmdUpdateScaFeeds, cmdUpdateProject CmdUpdateProject) *CmdUpdate {
	cmd := &cobra.Command{
		Use:               "update",
		Short:             "Update resources",
		Long:              `Update project sources, SBOM files, SCA feeds, and settings on the server.`,
		PersistentPreRunE: _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdUpdateSources.Command)
	cmd.AddCommand(cmdUpdateSbom.Command)
	cmd.AddCommand(cmdUpdateScaFeeds.Command)
	cmd.AddCommand(cmdUpdateProject.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdUpdate{cmd}
}
