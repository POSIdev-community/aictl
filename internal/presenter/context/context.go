package context

import (
	"github.com/spf13/cobra"
)

type CmdContext struct {
	*cobra.Command
}

func NewContextCmd(
	cmdConfigClear CmdConfigClear,
	cmdConfigSet CmdConfigSet,
	cmdConfigShow CmdConfigShow,
	cmdConfigUnset CmdConfigUnset) *CmdContext {

	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "Manage local aictl context",
		Long:  `Read and write the local aictl context (URI, token, TLS skip, project/branch ids) stored under ~/.config/aictl/.`,
	}

	cmd.AddCommand(cmdConfigClear.Command)
	cmd.AddCommand(cmdConfigSet.Command)
	cmd.AddCommand(cmdConfigShow.Command)
	cmd.AddCommand(cmdConfigUnset.Command)

	return &CmdContext{cmd}
}
