package application

import "github.com/spf13/cobra"

// revealDeprecatedCommandsForDocs temporarily clears Deprecated so cobra's
// GenMarkdownTree still emits pages and SEE ALSO links. Shell completion and
// help continue to hide those commands via Command.Deprecated at runtime.
func revealDeprecatedCommandsForDocs(root *cobra.Command) (restore func()) {
	type saved struct {
		cmd        *cobra.Command
		deprecated string
	}

	var list []saved

	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if cmd.Deprecated != "" {
			list = append(list, saved{cmd: cmd, deprecated: cmd.Deprecated})
			cmd.Deprecated = ""
		}
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)

	return func() {
		for _, item := range list {
			item.cmd.Deprecated = item.deprecated
		}
	}
}
