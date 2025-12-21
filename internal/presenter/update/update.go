package update

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdUpdate struct {
	*cobra.Command
}

var (
	projectIdFlag string
	branchIdFlag  string
)

func NewUpdateCmd(cfg *config.Config, cmdUpdateSources CmdUpdateSources) *CmdUpdate {
	cmd := &cobra.Command{
		Use:               "update",
		Short:             "Update resources",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdUpdateSources.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdUpdate{cmd}
}
