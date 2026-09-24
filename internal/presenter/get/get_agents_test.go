package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeAgentsUC struct {
	called int
	quiet  bool
}

func (f *fakeAgentsUC) Execute(_ context.Context, quiet bool) error {
	f.called++
	f.quiet = quiet
	return nil
}

func TestGetAgentsCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	cfg := mustGetCfg(t)
	uc := &fakeAgentsUC{}
	ucs := defaultGetUCs()
	ucs.agents = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "agents", "-q"))
	require.Equal(t, 1, uc.called)
	require.True(t, uc.quiet)

	t.Run("default_not_quiet", func(t *testing.T) {
		resetGetPackageFlags()
		uc2 := &fakeAgentsUC{}
		ucs2 := defaultGetUCs()
		ucs2.agents = uc2
		root2 := buildGetRoot(t, mustGetCfg(t), ucs2)
		require.NoError(t, cmdtest.Execute(t, root2.Command, "agents"))
		require.False(t, uc2.quiet)
	})
}
