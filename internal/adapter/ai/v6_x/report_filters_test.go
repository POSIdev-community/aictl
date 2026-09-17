package v6_x

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestToUserReportFiltersModel6x(t *testing.T) {
	t.Parallel()

	tTrue := true
	f := report.Filters{
		Apply:       true,
		LevelHigh:   &tTrue,
		Suspected:   &tTrue,
		Types:       []string{},
		Languages:   []string{"Dart"},
		ScanModules: []string{"SecretDetection", "BlackBox"},
	}

	model, err := toUserReportFiltersModel6x(f)
	require.NoError(t, err)
	require.NotNil(t, model.LevelHigh)
	require.NotNil(t, model.Suspected)
	require.Nil(t, model.LevelLow)
	require.Nil(t, model.Limit)
	require.Nil(t, model.NoPlaceToFix)
	require.Empty(t, *model.Types)
	require.Len(t, *model.Languages, 1)
	require.Len(t, *model.ScanModules, 2)
}
