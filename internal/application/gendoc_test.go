package application

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestRevealDeprecatedCommandsForDocs(t *testing.T) {
	child := &cobra.Command{
		Use:        "start",
		Short:      "deprecated alias",
		Deprecated: "use other",
		Run:        func(cmd *cobra.Command, args []string) {},
	}
	root := &cobra.Command{Use: "scan"}
	root.AddCommand(child)

	require.False(t, child.IsAvailableCommand())

	restore := revealDeprecatedCommandsForDocs(root)
	require.Empty(t, child.Deprecated)
	require.True(t, child.IsAvailableCommand())

	restore()
	require.Equal(t, "use other", child.Deprecated)
	require.False(t, child.IsAvailableCommand())
}
