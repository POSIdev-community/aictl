package scan

import (
	"fmt"
	"strings"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/spf13/cobra"
)

type PersistentPreRunEScanStartCmd _utils.RunE

type CmdScanStart struct {
	*cobra.Command
}

var scanLabel string

func NewPersistentPreRunEScanStartCmd(prev PersistentPreRunEScanCmd) PersistentPreRunEScanStartCmd {
	return _utils.ChainRunE(prev, func(cmd *cobra.Command, args []string) error {
		if scanLabel != "" {
			if len(scanLabel) > 40 {
				return fmt.Errorf("label length must be less than 40")
			}

			if strings.ContainsAny(scanLabel, "#%?/;,\"\r\n\\") {
				return fmt.Errorf("label contains invalid characters")
			}
		}

		return nil
	})
}

func NewScanStartCmd(persistentPreRunE PersistentPreRunEScanStartCmd, cmdScanStart CmdScanStartBranch,
	cmdScanStartProject CmdScanStartProject) CmdScanStart {
	cmd := &cobra.Command{
		Use:               "start",
		Short:             "Start scan",
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(cmdScanStart.Command)
	cmd.AddCommand(cmdScanStartProject.Command)

	cmd.PersistentFlags().StringVar(&scanLabel, "scan-label", "", "scan label for scan")

	return CmdScanStart{cmd}
}
