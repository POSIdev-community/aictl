package context

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeConfigUnsetUC struct {
	called                               int
	uri, token, tls, projectID, branchID bool
}

func (f *fakeConfigUnsetUC) Execute(uriUnset, tokenUnset, tlsUnset, projectIdUnset, branchIdUnset bool) error {
	f.called++
	f.uri, f.token, f.tls = uriUnset, tokenUnset, tlsUnset
	f.projectID, f.branchID = projectIdUnset, branchIdUnset
	return nil
}

func TestConfigUnsetCommand(t *testing.T) {
	t.Run("requires_at_least_one_flag", func(t *testing.T) {
		uc := &fakeConfigUnsetUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, noopShowUC{}, uc)
		require.Error(t, cmdtest.Execute(t, root.Command, "unset"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("passes_all_flags", func(t *testing.T) {
		uc := &fakeConfigUnsetUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, noopShowUC{}, uc)
		require.NoError(t, cmdtest.Execute(t, root.Command, "unset", "-u", "-t", "-p", "-b", "--tls-skip"))
		require.Equal(t, 1, uc.called)
		require.True(t, uc.uri && uc.token && uc.tls && uc.projectID && uc.branchID)
	})
}
