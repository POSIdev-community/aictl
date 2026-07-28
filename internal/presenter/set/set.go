package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type PersistentPreRunESetCmd _utils.RunE

func NewPersistentPreRunESetCmd(cfg *config.Config) PersistentPreRunESetCmd {
	return _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg))
}

type CmdSet struct {
	*cobra.Command
}

func NewSetCmd(
	persistentPreRunESetCmd PersistentPreRunESetCmd,
	setProjectCmd CmdSetProject) *CmdSet {
	cmd := &cobra.Command{
		Use:               "set",
		Short:             "Set resource configuration",
		Long:              `Replace or upload project configuration on the server (settings, policies, exclusions).`,
		PersistentPreRunE: persistentPreRunESetCmd,
	}

	cmd.AddCommand(setProjectCmd.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdSet{cmd}
}
