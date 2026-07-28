package get

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
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
	cmdGetHealthcheck CmdGetHealthcheck,
	cmdGetProjects CmdGetProjects,
	cmdGetProject CmdGetProject,
	cmdGetBranches CmdGetBranches,
	cmdGetBranch CmdGetBranch,
	cmdGetScans CmdGetScans,
	cmdGetScan CmdGetScan,
	cmdGetAgents CmdGetAgents,
	cmdGetVersion CmdGetVersion) *CmdGet {

	cmd := &cobra.Command{
		Use:               "get",
		Short:             "Get resources",
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdGetHealthcheck.Command)
	cmd.AddCommand(cmdGetProjects.Command)
	cmd.AddCommand(cmdGetProject.Command)
	cmd.AddCommand(cmdGetBranches.Command)
	cmd.AddCommand(cmdGetBranch.Command)
	cmd.AddCommand(cmdGetScans.Command)
	cmd.AddCommand(cmdGetScan.Command)
	cmd.AddCommand(cmdGetAgents.Command)
	cmd.AddCommand(cmdGetVersion.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdGet{cmd}
}
