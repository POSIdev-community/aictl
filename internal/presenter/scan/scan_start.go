package scan

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunEScanStartCmd _utils.RunE

type CmdScanStart struct {
	*cobra.Command
}

var (
	scanLabel string
	fullScan  bool
)

func scanTypeFromFlags() scantype.Type {
	if fullScan {
		return scantype.Full
	}

	return scantype.Incremental
}

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
	cmd.PersistentFlags().BoolVar(&fullScan, "full-scan", false, "run full scan instead of incremental")

	return CmdScanStart{cmd}
}
