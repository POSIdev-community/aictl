package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeHealthcheckUC struct{ called int }

func (f *fakeHealthcheckUC) Execute(context.Context) error { f.called++; return nil }

func TestGetHealthcheckCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	uc := &fakeHealthcheckUC{}
	ucs := defaultGetUCs()
	ucs.healthcheck = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "healthcheck"))
	require.Equal(t, 1, uc.called)
}
