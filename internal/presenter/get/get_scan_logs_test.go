package get

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeScanLogsUC struct {
	called int
	scanID uuid.UUID
	out    string
}

func (f *fakeScanLogsUC) Execute(_ context.Context, scanId uuid.UUID, outputPath string) error {
	f.called++
	f.scanID, f.out = scanId, outputPath
	return nil
}

func TestGetScanLogsCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	out := filepath.Join(t.TempDir(), "logs.txt")
	uc := &fakeScanLogsUC{}
	ucs := defaultGetUCs()
	ucs.scanLogs = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "scan", "logs", scanID.String(), "-p", projectID.String(), "-o", out))
	require.Equal(t, scanID, uc.scanID)
	require.Equal(t, out, uc.out)

	t.Run("force_overwrites", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "logs.txt")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc2 := &fakeScanLogsUC{}
		ucs2 := defaultGetUCs()
		ucs2.scanLogs = uc2
		root2 := buildGetRoot(t, mustGetCfg(t), ucs2)
		require.NoError(t, cmdtest.Execute(t, root2.Command, "scan", "logs", scanID.String(), "-p", projectID.String(), "-o", exist, "-f"))
		require.Equal(t, exist, uc2.out)
	})
}
