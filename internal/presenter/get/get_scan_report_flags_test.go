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

type fakeCustomReportUC struct {
	called                                       int
	scanID                                       uuid.UUID
	name, out, l10n                              string
	includeComments, includeDFD, includeGlossary bool
}

func (f *fakeCustomReportUC) Execute(
	_ context.Context,
	scanId uuid.UUID,
	customReportName, outPath string,
	includeComments, includeDFD, includeGlossary bool,
	l10n string,
) error {
	f.called++
	f.scanID = scanId
	f.name, f.out, f.l10n = customReportName, outPath, l10n
	f.includeComments, f.includeDFD, f.includeGlossary = includeComments, includeDFD, includeGlossary
	return nil
}

func TestGetScanReportFlags(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	t.Run("all_report_persistent_flags", func(t *testing.T) {
		resetGetPackageFlags()
		out := filepath.Join(t.TempDir(), "report.html")
		uc := &fakeCustomReportUC{}
		ucs := defaultGetUCs()
		ucs.customReport = uc
		root := buildGetRoot(t, mustGetCfg(t), ucs)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"scan", "report", "MyTemplate", scanID.String(),
			"-p", projectID.String(),
			"-o", out,
			"--include-comments",
			"--include-dfd",
			"--include-glossary",
			"--localization", "ru",
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, scanID, uc.scanID)
		require.Equal(t, "MyTemplate", uc.name)
		require.Equal(t, out, uc.out)
		require.True(t, uc.includeComments && uc.includeDFD && uc.includeGlossary)
		require.Equal(t, "ru", uc.l10n)
	})

	t.Run("force_required_when_output_exists", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "exists.html")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc := &fakeCustomReportUC{}
		ucs := defaultGetUCs()
		ucs.customReport = uc
		root := buildGetRoot(t, mustGetCfg(t), ucs)
		require.Error(t, cmdtest.Execute(t, root.Command,
			"scan", "report", "MyTemplate", scanID.String(),
			"-p", projectID.String(),
			"-o", exist,
		))
		require.Equal(t, 0, uc.called)
	})

	t.Run("force_overwrites", func(t *testing.T) {
		resetGetPackageFlags()
		exist := filepath.Join(t.TempDir(), "exists.html")
		require.NoError(t, os.WriteFile(exist, []byte("x"), 0o644))
		uc := &fakeCustomReportUC{}
		ucs := defaultGetUCs()
		ucs.customReport = uc
		root := buildGetRoot(t, mustGetCfg(t), ucs)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"scan", "report", "MyTemplate", scanID.String(),
			"-p", projectID.String(),
			"-o", exist,
			"-f",
		))
		require.Equal(t, exist, uc.out)
	})
}
