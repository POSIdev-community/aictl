package delete

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdDelete struct {
	*cobra.Command
}

func NewDeleteCmd(cfg *config.Config, cmdDeleteProjects CmdDeleteProjects) *CmdDelete {
	cmd := &cobra.Command{
		Use:               "delete",
		Short:             "Delete resources",
		PersistentPreRunE: _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdDeleteProjects.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdDelete{cmd}
}
