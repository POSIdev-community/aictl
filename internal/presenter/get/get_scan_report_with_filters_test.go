package get

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

func TestGetScanReportWithFilters_RequiresFilterFlag(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()

	cfg := mustGetCfg(t)
	root := buildGetRoot(t, cfg, defaultGetUCs())

	scanID := uuid.New().String()
	err := cmdtest.Execute(t, root.Command, "scan", "report", "with-filters", "sarif", scanID, "-p", uuid.New().String())
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one filter flag is required")
}

func TestGetScanReportWithFilters_AcceptsLevelHigh(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()

	cfg := mustGetCfg(t)
	root := buildGetRoot(t, cfg, defaultGetUCs())

	scanID := uuid.New().String()
	err := cmdtest.Execute(t, root.Command, "scan", "report", "with-filters", "sarif", scanID, "-p", uuid.New().String(), "--level-high")
	require.NoError(t, err)
	require.NotNil(t, reportFilters.LevelHigh)
	require.True(t, *reportFilters.LevelHigh)
	require.Nil(t, reportFilters.LevelLow)
	require.Empty(t, reportFilters.Types)
	require.Empty(t, reportFilters.Languages)
	require.Empty(t, reportFilters.ScanModules)
}

func TestGetScanReportWithFilters_RejectsEmptyType(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()

	cfg := mustGetCfg(t)
	root := buildGetRoot(t, cfg, defaultGetUCs())

	scanID := uuid.New().String()
	err := cmdtest.Execute(t, root.Command, "scan", "report", "with-filters", "sarif", scanID, "-p", uuid.New().String(), "--type", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "type")
}

func TestGetScanReportWithFilters_RejectsLanguageNone(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()

	cfg := mustGetCfg(t)
	root := buildGetRoot(t, cfg, defaultGetUCs())

	scanID := uuid.New().String()
	err := cmdtest.Execute(t, root.Command, "scan", "report", "with-filters", "sarif", scanID, "-p", uuid.New().String(), "--language", "None")
	require.Error(t, err)
	require.Contains(t, err.Error(), "None")
}

func TestGetScanReportWithFilters_RejectsForbiddenScanModule(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()

	cfg := mustGetCfg(t)
	root := buildGetRoot(t, cfg, defaultGetUCs())

	scanID := uuid.New().String()
	err := cmdtest.Execute(t, root.Command, "scan", "report", "with-filters", "sarif", scanID, "-p", uuid.New().String(), "--scan-module", "DataFlowAnalysis")
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan-module")
}
