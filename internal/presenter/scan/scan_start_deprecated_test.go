package scan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanStartCommandsAreDeprecated(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()

	cfg := mustScanCfg(t)
	root := buildScanRoot(t, cfg, noopAwaitUC{}, noopCheckPoliciesUC{}, noopStartBranchUC{}, noopStartProjectUC{}, noopStopUC{})

	start, _, err := root.Find([]string{"start"})
	require.NoError(t, err)
	require.NotEmpty(t, start.Deprecated)
	require.False(t, start.IsAvailableCommand(), "deprecated start must be hidden from completion/help")

	branch, _, err := root.Find([]string{"start", "branch"})
	require.NoError(t, err)
	require.NotEmpty(t, branch.Deprecated)
	require.False(t, branch.IsAvailableCommand())

	project, _, err := root.Find([]string{"start", "project"})
	require.NoError(t, err)
	require.NotEmpty(t, project.Deprecated)
	require.False(t, project.IsAvailableCommand())
}
