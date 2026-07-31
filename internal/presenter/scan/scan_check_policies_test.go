package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeCheckPoliciesUC struct {
	called                 int
	scanID                 uuid.UUID
	failOnPoliciesRejected bool
}

func (f *fakeCheckPoliciesUC) Execute(_ context.Context, scanId uuid.UUID, failOnPoliciesRejected bool) error {
	f.called++
	f.scanID, f.failOnPoliciesRejected = scanId, failOnPoliciesRejected
	return nil
}

func TestScanCheckPoliciesCmd(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	cfg := mustScanCfg(t)
	uc := &fakeCheckPoliciesUC{}
	root := buildScanRoot(t, cfg, noopAwaitUC{}, uc, noopStartBranchUC{}, noopStartProjectUC{}, noopStopUC{})
	require.NoError(t, cmdtest.Execute(t, root.Command,
		"check-policies", scanID.String(),
		"-p", projectID.String(),
		"--fail-on-policies-rejected",
	))
	require.Equal(t, scanID, uc.scanID)
	require.True(t, uc.failOnPoliciesRejected)
}
