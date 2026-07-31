package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeVersionUC struct{ called int }

func (f *fakeVersionUC) Execute(context.Context) error { f.called++; return nil }

func TestGetVersionCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	uc := &fakeVersionUC{}
	ucs := defaultGetUCs()
	ucs.version = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "version"))
	require.Equal(t, 1, uc.called)
}
