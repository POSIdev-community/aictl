package scan

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type CmdScan struct {
	*cobra.Command
}

func NewScanCmd(
	cfg *config.Config,
	cmdScanAwait CmdScanAwait,
	cmdScanStart CmdScanStart,
	cmdScanStop CmdScanStop) *CmdScan {

	cmd := &cobra.Command{
		Use:               "scan",
		Short:             "Scan ",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdScanAwait.Command)
	cmd.AddCommand(cmdScanStart.Command)
	cmd.AddCommand(cmdScanStop.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdScan{cmd}
}
