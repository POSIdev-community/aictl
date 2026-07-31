package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScansUC struct {
	called int
	quite  bool
	latest bool
	filter regexfilter.RegexFilter
}

func (f *fakeScansUC) Execute(_ context.Context, filter regexfilter.RegexFilter, quite, latest bool) error {
	f.called++
	f.filter, f.quite, f.latest = filter, quite, latest
	return nil
}

func TestGetScansCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	branchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	cfg := mustGetCfg(t)
	uc := &fakeScansUC{}
	ucs := defaultGetUCs()
	ucs.scans = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scans", "night", "-q", "--latest", "-b", branchID.String()))
	require.True(t, uc.quite && uc.latest)
	require.Equal(t, branchID, cfg.BranchId())
}
