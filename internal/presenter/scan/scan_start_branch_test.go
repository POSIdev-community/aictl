package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeStartBranchUC struct {
	called    int
	scanLabel string
	scanType  scantype.Type
}

func (f *fakeStartBranchUC) Execute(_ context.Context, scanLabel string, scanType scantype.Type) error {
	f.called++
	f.scanLabel, f.scanType = scanLabel, scanType
	return nil
}

func TestScanStartBranchCmd(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	branchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	cfg := mustScanCfg(t)
	uc := &fakeStartBranchUC{}
	root := buildScanRoot(t, cfg, noopAwaitUC{}, noopCheckPoliciesUC{}, uc, noopStartProjectUC{}, noopStopUC{})
	require.NoError(t, cmdtest.Execute(t, root.Command,
		"start", "branch", branchID.String(),
		"-p", projectID.String(),
		"--scan-label", "nightly",
		"--full-scan",
	))
	require.Equal(t, "nightly", uc.scanLabel)
	require.Equal(t, scantype.Full, uc.scanType)
	require.Equal(t, projectID, cfg.ProjectId())
	require.Equal(t, branchID, cfg.BranchId())
}
