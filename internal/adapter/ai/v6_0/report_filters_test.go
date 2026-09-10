package v6_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestToUserReportFiltersModel60(t *testing.T) {
	t.Parallel()

	tTrue := true
	f := report.Filters{
		Apply:       true,
		LevelHigh:   &tTrue,
		Types:       []string{},
		Languages:   []string{"OneC"},
		ScanModules: []string{"SecretDetection", "MaliciousCodeDetection"},
	}

	model, err := toUserReportFiltersModel60(f)
	require.NoError(t, err)
	require.NotNil(t, model.LevelHigh)
	require.Len(t, *model.Languages, 1)
	require.Len(t, *model.ScanModules, 2)
	require.Nil(t, model.Limit)
	require.Nil(t, model.NoPlaceToFix)
}

func TestMapProgrammingLanguages60_RejectsDart(t *testing.T) {
	t.Parallel()
	_, err := mapProgrammingLanguages60([]string{"Dart"})
	require.Error(t, err)
}
