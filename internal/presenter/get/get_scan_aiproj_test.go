package get

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScanAiprojUC struct {
	called int
	scanID uuid.UUID
	out    string
}

func (f *fakeScanAiprojUC) Execute(_ context.Context, scanId uuid.UUID, outputPath string) error {
	f.called++
	f.scanID, f.out = scanId, outputPath
	return nil
}

func TestGetScanAiprojCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	cfg := mustGetCfg(t)
	out := filepath.Join(t.TempDir(), "scan.aiproj")
	uc := &fakeScanAiprojUC{}
	ucs := defaultGetUCs()
	ucs.scanAiproj = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scan", "aiproj", scanID.String(), "-p", projectID.String(), "-o", out))
	require.Equal(t, scanID, uc.scanID)
	require.Equal(t, out, uc.out)
}
