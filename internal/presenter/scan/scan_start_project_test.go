package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeStartProjectUC struct {
	called    int
	scanLabel string
	scanType  scantype.Type
}

func (f *fakeStartProjectUC) Execute(_ context.Context, scanLabel string, scanType scantype.Type) error {
	f.called++
	f.scanLabel, f.scanType = scanLabel, scanType
	return nil
}

func TestScanStartProjectCmd(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustScanCfg(t)
	uc := &fakeStartProjectUC{}
	root := buildScanRoot(t, cfg, noopAwaitUC{}, noopCheckPoliciesUC{}, noopStartBranchUC{}, uc, noopStopUC{})
	require.NoError(t, cmdtest.Execute(t, root.Command,
		"start", "project", projectID.String(),
		"--scan-label", "release",
		"--full-scan",
	))
	require.Equal(t, "release", uc.scanLabel)
	require.Equal(t, scantype.Full, uc.scanType)
	require.Equal(t, projectID, cfg.ProjectId())
}
