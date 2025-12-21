package set

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdSet struct {
	*cobra.Command
}

func NewSetCmd(
	cfg *config.Config,
	setProjectCmd CmdSetProject) *CmdSet {
	cmd := &cobra.Command{
		Use:               "set",
		Short:             "Set",
		PersistentPreRunE: _utils.ConcatFuncs(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(setProjectCmd.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	return &CmdSet{cmd}
}
