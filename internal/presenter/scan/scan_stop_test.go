package scan

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeStopUC struct {
	called int
	scanID uuid.UUID
}

func (f *fakeStopUC) Execute(_ context.Context, scanResultId uuid.UUID) error {
	f.called++
	f.scanID = scanResultId
	return nil
}

func TestScanStopCmd(t *testing.T) {
	t.Cleanup(resetScanStartFlags)
	resetScanStartFlags()
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	cfg := mustScanCfg(t)
	uc := &fakeStopUC{}
	root := buildScanRoot(t, cfg, noopAwaitUC{}, noopCheckPoliciesUC{}, noopStartBranchUC{}, noopStartProjectUC{}, uc)
	require.NoError(t, cmdtest.Execute(t, root.Command, "stop", scanID.String()))
	require.Equal(t, scanID, uc.scanID)
}
