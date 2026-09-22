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

type fakeScanStatisticUC struct {
	called     int
	scanID     uuid.UUID
	out        string
	json       bool
	withTriage bool
}

func (f *fakeScanStatisticUC) Execute(_ context.Context, scanId uuid.UUID, outPath string, json, withTriage bool) error {
	f.called++
	f.scanID, f.out, f.json, f.withTriage = scanId, outPath, json, withTriage
	return nil
}

func TestGetScanStatisticCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	out := filepath.Join(t.TempDir(), "stat.json")
	uc := &fakeScanStatisticUC{}
	ucs := defaultGetUCs()
	ucs.scanStatistic = uc
	root := buildGetRoot(t, mustGetCfg(t), ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command,
		"scan", "statistic", scanID.String(),
		"-p", projectID.String(),
		"-o", out,
		"--json",
		"--with-triage",
	))
	require.Equal(t, scanID, uc.scanID)
	require.Equal(t, out, uc.out)
	require.True(t, uc.json)
	require.True(t, uc.withTriage)

	t.Run("force_overwrites", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "stat.json")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc2 := &fakeScanStatisticUC{}
		ucs2 := defaultGetUCs()
		ucs2.scanStatistic = uc2
		root2 := buildGetRoot(t, mustGetCfg(t), ucs2)
		require.NoError(t, cmdtest.Execute(t, root2.Command,
			"scan", "statistic", scanID.String(),
			"-p", projectID.String(),
			"-o", exist,
			"-f",
		))
		require.Equal(t, exist, uc2.out)
		require.False(t, uc2.withTriage)
	})
}
