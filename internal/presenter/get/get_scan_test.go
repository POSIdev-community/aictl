package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScanUC struct {
	called int
	scanID uuid.UUID
}

func (f *fakeScanUC) Execute(_ context.Context, id uuid.UUID) error {
	f.called++
	f.scanID = id
	return nil
}

func TestGetScanCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	uc := &fakeScanUC{}
	ucs := defaultGetUCs()
	ucs.scan = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scan", scanID.String(), "-p", projectID.String()))
	require.Equal(t, scanID, uc.scanID)
}
