package scan

import (
	"github.com/spf13/cobra"
)

type CmdScanStart struct {
	*cobra.Command
}

func NewScanStartCmd(cmdScanStart CmdScanStartBranch, cmdScanStartProject CmdScanStartProject) CmdScanStart {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start scan",
	}

	cmd.AddCommand(cmdScanStart.Command)
	cmd.AddCommand(cmdScanStartProject.Command)

	return CmdScanStart{cmd}
}
