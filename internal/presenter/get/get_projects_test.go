package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeProjectsUC struct {
	called int
	filter regexfilter.RegexFilter
	quite  bool
}

func (f *fakeProjectsUC) Execute(_ context.Context, filter regexfilter.RegexFilter, quite bool) error {
	f.called++
	f.filter, f.quite = filter, quite
	return nil
}

func TestGetProjectsCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	cfg := mustGetCfg(t)
	uc := &fakeProjectsUC{}
	ucs := defaultGetUCs()
	ucs.projects = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "projects", "my-.*", "-q"))
	require.Equal(t, 1, uc.called)
	require.True(t, uc.quite)
}
