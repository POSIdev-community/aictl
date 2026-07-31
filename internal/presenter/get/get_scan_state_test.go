package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScanStateUC struct {
	called           int
	scanID           uuid.UUID
	failOnScanFailed bool
}

func (f *fakeScanStateUC) Execute(_ context.Context, scanId uuid.UUID, failOnScanFailed bool) error {
	f.called++
	f.scanID, f.failOnScanFailed = scanId, failOnScanFailed
	return nil
}

func TestGetScanStateCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	uc := &fakeScanStateUC{}
	ucs := defaultGetUCs()
	ucs.scanState = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scan", "stage", scanID.String(), "-p", projectID.String(), "--fail-on-scan-failed"))
	require.Equal(t, scanID, uc.scanID)
	require.True(t, uc.failOnScanFailed)
}
