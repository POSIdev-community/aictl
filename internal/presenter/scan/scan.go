package scan

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type PersistentPreRunEScanCmd _utils.RunE

type CmdScan struct {
	*cobra.Command
}

func NewPersistentPreRunEScanCmd(cfg *config.Config) PersistentPreRunEScanCmd {
	return _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg))
}

func NewScanCmd(
	persistentPreRunE PersistentPreRunEScanCmd,
	cmdScanAwait CmdScanAwait,
	cmdScanStart CmdScanStart,
	cmdScanStop CmdScanStop) *CmdScan {

	cmd := &cobra.Command{
		Use:               "scan",
		Short:             "Scan ",
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdScanAwait.Command)
	cmd.AddCommand(cmdScanStart.Command)
	cmd.AddCommand(cmdScanStop.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdScan{cmd}
}
