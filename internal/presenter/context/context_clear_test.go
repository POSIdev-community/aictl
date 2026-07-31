package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeConfigClearUC struct {
	called      int
	skipConfirm bool
}

func (f *fakeConfigClearUC) Execute(_ context.Context, skipConfirm bool) error {
	f.called++
	f.skipConfirm = skipConfirm
	return nil
}

func TestConfigClearCommand(t *testing.T) {
	t.Run("with_yes", func(t *testing.T) {
		uc := &fakeConfigClearUC{}
		root := newCtxRoot(uc, noopSetUC{}, noopShowUC{}, noopUnsetUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "clear", "-y"))
		require.Equal(t, 1, uc.called)
		require.True(t, uc.skipConfirm)
	})

	t.Run("without_yes", func(t *testing.T) {
		uc := &fakeConfigClearUC{}
		root := newCtxRoot(uc, noopSetUC{}, noopShowUC{}, noopUnsetUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "clear"))
		require.Equal(t, 1, uc.called)
		require.False(t, uc.skipConfirm)
	})
}
