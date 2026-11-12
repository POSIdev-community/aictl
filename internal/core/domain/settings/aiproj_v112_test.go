package settings

import (
	"testing"

	"github.com/stretchr/testify/require"

	v112 "github.com/POSIdev-community/aiproj/model/v1_12"
)

func TestUpdateFromV112_TriageAndRaceAsint(t *testing.T) {
	t.Parallel()

	enabled := true
	severity := v112.AIProjAutomatedAiTriageMinSeverityHigh
	s := &ScanSettings{}
	err := s.UpdateFromV112(&v112.AIProj{
		ProjectName:                  "demo",
		ProgrammingLanguages:         []v112.AIProjProgrammingLanguagesElem{v112.AIProjProgrammingLanguagesElemGo},
		ScanModules:                  []v112.AIProjScanModulesElem{v112.AIProjScanModulesElemRaceConditionDetection, v112.AIProjScanModulesElemAsint},
		AutomatedAiTriage:            &enabled,
		AutomatedAiTriageMinSeverity: &severity,
	})
	require.NoError(t, err)
	require.Equal(t, "demo", s.ProjectName)
	require.True(t, s.WhiteBoxSettings.SearchForRaceConditionsEnabled)
	require.True(t, s.WhiteBoxSettings.AsintEnabled)
	require.True(t, s.AutomatedAiTriageEnabled)
	require.Equal(t, "High", s.AutomatedAiTriageMinSeverity)
}
