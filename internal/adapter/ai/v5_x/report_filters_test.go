package v5_x

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestToUserReportFiltersModel5x(t *testing.T) {
	t.Parallel()

	tTrue := true
	f := report.Filters{
		Apply:     true,
		LevelHigh: &tTrue,
		Types:     []string{},
		Languages: []string{"Java"},
		ScanModules: []string{
			"StaticCodeAnalysis",
			"BlackBox",
		},
	}

	model, err := toUserReportFiltersModel5x(f)
	require.NoError(t, err)
	require.NotNil(t, model.LevelHigh)
	require.True(t, *model.LevelHigh)
	require.Nil(t, model.LevelLow)
	require.Nil(t, model.Limit)
	require.Nil(t, model.NoPlaceToFix)
	require.NotNil(t, model.Types)
	require.Empty(t, *model.Types)
	require.Len(t, *model.Languages, 1)
	require.Len(t, *model.ScanModules, 2)
}

func TestMapScanModules5x_RejectsSecretDetection(t *testing.T) {
	t.Parallel()
	_, err := mapScanModules5x([]string{"SecretDetection"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan-module")
}
