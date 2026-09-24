package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeBranchesUC struct {
	called int
	quiet  bool
	filter regexfilter.RegexFilter
}

func (f *fakeBranchesUC) Execute(_ context.Context, filter regexfilter.RegexFilter, quiet bool) error {
	f.called++
	f.filter, f.quiet = filter, quiet
	return nil
}

func TestGetBranchesCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustGetCfg(t)
	uc := &fakeBranchesUC{}
	ucs := defaultGetUCs()
	ucs.branches = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "branches", "main", "-q", "-p", projectID.String()))
	require.True(t, uc.quiet)
	require.Equal(t, projectID, cfg.ProjectId())
}
