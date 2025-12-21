package get

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdGet struct {
	*cobra.Command
}

func NewGetCmd(
	cfg *config.Config,
	cmdGetProjects CmdGetProjects,
	cmdGetReports CmdGetReports,
	cmdGetScan CmdGetScan,
	cmdGetScans CmdGetScans) *CmdGet {

	cmd := &cobra.Command{
		Use:               "get",
		Short:             "Get resources",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdGetProjects.Command)
	cmd.AddCommand(cmdGetReports.Command)
	cmd.AddCommand(cmdGetScan.Command)
	cmd.AddCommand(cmdGetScans.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdGet{cmd}
}
