package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeAwaitUC struct {
	called           int
	scanID           uuid.UUID
	failOnScanFailed bool
}

func (f *fakeAwaitUC) Execute(_ context.Context, scanId uuid.UUID, failOnScanFailed bool) error {
	f.called++
	f.scanID, f.failOnScanFailed = scanId, failOnScanFailed
	return nil
}

func TestScanAwaitCmd(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	cfg := mustScanCfg(t)
	uc := &fakeAwaitUC{}
	root := buildScanRoot(t, cfg, uc, noopCheckPoliciesUC{}, noopStartBranchUC{}, noopStartProjectUC{}, noopStopUC{})
	require.NoError(t, cmdtest.Execute(t, root.Command, "await", scanID.String(), "-p", projectID.String(), "--fail-on-scan-failed"))
	require.Equal(t, scanID, uc.scanID)
	require.True(t, uc.failOnScanFailed)
	require.Equal(t, projectID, cfg.ProjectId())
}
