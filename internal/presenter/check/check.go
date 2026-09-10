package check

import (
	"github.com/spf13/cobra"

	_utils "github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunECheckCmd _utils.RunE

func NewPersistentPreRunECheckCmd() PersistentPreRunECheckCmd {
	return _utils.InitializeLogger
}

type CmdCheck struct {
	*cobra.Command
}

func NewCheckCmd(persistentPreRunE PersistentPreRunECheckCmd, aiprojCmd CmdCheckAiproj) *CmdCheck {
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "Validate local resources",
		Long:              `Offline checks for local files and configuration (no AI server connection required).`,
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(aiprojCmd.Command)
	_utils.AddVerbosePersistentFlags(cmd)

	return &CmdCheck{cmd}
}
