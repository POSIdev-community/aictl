package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeConfigShowUC struct {
	called     int
	json, yaml bool
}

func (f *fakeConfigShowUC) Execute(_ context.Context, json, yaml bool) error {
	f.called++
	f.json, f.yaml = json, yaml
	return nil
}

func TestConfigShowCommand(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		uc := &fakeConfigShowUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, uc, noopUnsetUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "show"))
		require.Equal(t, 1, uc.called)
		require.False(t, uc.json)
		require.False(t, uc.yaml)
	})

	t.Run("json", func(t *testing.T) {
		uc := &fakeConfigShowUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, uc, noopUnsetUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "show", "--json"))
		require.Equal(t, 1, uc.called)
		require.True(t, uc.json)
		require.False(t, uc.yaml)
	})

	t.Run("yaml", func(t *testing.T) {
		uc := &fakeConfigShowUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, uc, noopUnsetUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "show", "--yaml"))
		require.True(t, uc.yaml)
		require.False(t, uc.json)
	})

	t.Run("both_rejected", func(t *testing.T) {
		uc := &fakeConfigShowUC{}
		root := newCtxRoot(noopClearUC{}, noopSetUC{}, uc, noopUnsetUC{})
		require.Error(t, cmdtest.Execute(t, root.Command, "show", "--json", "--yaml"))
		require.Equal(t, 0, uc.called)
	})
}
