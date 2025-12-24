package get

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type PersistentPreRunEGetCmd _utils.RunE

type CmdGet struct {
	*cobra.Command
}

func NewPersistentPreRunEGetCmd(cfg *config.Config) PersistentPreRunEGetCmd {
	return _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg))
}

func NewGetCmd(
	persistentPreRunE PersistentPreRunEGetCmd,
	cmdGetProjects CmdGetProjects,
	cmdGetScan CmdGetScan) *CmdGet {

	cmd := &cobra.Command{
		Use:               "get",
		Short:             "Get resources",
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdGetProjects.Command)
	cmd.AddCommand(cmdGetScan.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdGet{cmd}
}
