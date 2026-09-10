package report_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestFilters_HasAny(t *testing.T) {
	t.Parallel()

	require.False(t, report.EmptyFilters().HasAny())

	tTrue := true
	require.True(t, report.Filters{LevelHigh: &tTrue}.HasAny())
	require.True(t, report.Filters{Types: []string{"XSS"}}.HasAny())
	require.True(t, report.Filters{Languages: []string{"Java"}}.HasAny())
	require.True(t, report.Filters{ScanModules: []string{"BlackBox"}}.HasAny())
}

func TestFilters_Normalize_Dedupes(t *testing.T) {
	t.Parallel()

	f := report.Filters{
		Types:       []string{"a", "b", "a"},
		Languages:   []string{"Java", "Java", "Go"},
		ScanModules: []string{"BlackBox", "BlackBox"},
	}
	f.Normalize()
	require.Equal(t, []string{"a", "b"}, f.Types)
	require.Equal(t, []string{"Java", "Go"}, f.Languages)
	require.Equal(t, []string{"BlackBox"}, f.ScanModules)
}

func TestFilters_ValidateArrays(t *testing.T) {
	t.Parallel()

	require.Error(t, report.Filters{Types: []string{""}}.ValidateArrays())
	require.Error(t, report.Filters{Languages: []string{"None"}}.ValidateArrays())
	require.Error(t, report.Filters{Languages: []string{"Haskell"}}.ValidateArrays())
	require.Error(t, report.Filters{ScanModules: []string{"DataFlowAnalysis"}}.ValidateArrays())
	require.Error(t, report.Filters{ScanModules: []string{"VulnerableSourceCode"}}.ValidateArrays())
	require.NoError(t, report.Filters{
		Types:       []string{"XSS"},
		Languages:   []string{"Java", "Dart"},
		ScanModules: []string{"StaticCodeAnalysis", "SecretDetection"},
	}.ValidateArrays())
}
