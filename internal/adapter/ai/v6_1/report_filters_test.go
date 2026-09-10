package v6_1

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestToUserReportFiltersModel61(t *testing.T) {
	t.Parallel()

	tTrue := true
	f := report.Filters{
		Apply:       true,
		LevelHigh:   &tTrue,
		Types:       []string{},
		Languages:   []string{"Dart", "OneC"},
		ScanModules: []string{"SecretDetection"},
	}

	model, err := toUserReportFiltersModel61(f)
	require.NoError(t, err)
	require.NotNil(t, model.LevelHigh)
	require.Len(t, *model.Languages, 2)
	require.Len(t, *model.ScanModules, 1)
	require.Nil(t, model.Limit)
	require.Nil(t, model.NoPlaceToFix)
}
