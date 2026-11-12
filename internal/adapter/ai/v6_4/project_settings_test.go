package v6_4

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_4"
)

func TestMapProjectSettingsFromModel_AsintAndAiTriage(t *testing.T) {
	t.Parallel()

	useNaira := true
	triageEnabled := true
	minSeverity := v6_4.IssueLevelMedium
	asintEnabled := true
	raceEnabled := true

	got := mapProjectSettingsFromModel(v6_4.ProjectSettingsModel{
		WhiteBoxSettings: &v6_4.WhiteBoxSettingsModel{
			AsintEnabled:                   &asintEnabled,
			SearchForRaceConditionsEnabled: &raceEnabled,
		},
		AsintSettings: &v6_4.AsintSettingsModel{
			UseNaira: &useNaira,
		},
		AutomatedAiTriageEnabled:     &triageEnabled,
		AutomatedAiTriageMinSeverity: &minSeverity,
	})

	require.True(t, got.WhiteBoxSettings.AsintEnabled)
	require.True(t, got.WhiteBoxSettings.SearchForRaceConditionsEnabled)
	require.True(t, got.AsintSettings.UseNaira)
	require.True(t, got.AutomatedAiTriageEnabled)
	require.Equal(t, string(v6_4.IssueLevelMedium), got.AutomatedAiTriageMinSeverity)
}

func TestMapProjectSettingsFromModel_AsintAndAiTriageDefaults(t *testing.T) {
	t.Parallel()

	got := mapProjectSettingsFromModel(v6_4.ProjectSettingsModel{})
	require.False(t, got.AsintSettings.UseNaira)
	require.False(t, got.AutomatedAiTriageEnabled)
	require.Empty(t, got.AutomatedAiTriageMinSeverity)
}

func TestMapProjectSettingsToModel_AsintAndAiTriage(t *testing.T) {
	t.Parallel()

	got := mapProjectSettingsToModel(&settings.ScanSettings{
		ProjectName: "demo",
		WhiteBoxSettings: settings.WhiteBoxSettings{
			SearchForRaceConditionsEnabled: true,
			AsintEnabled:                   true,
		},
		AsintSettings: settings.AsintSettings{
			UseNaira: true,
		},
		AutomatedAiTriageEnabled:     true,
		AutomatedAiTriageMinSeverity: "High",
	})

	require.NotNil(t, got.WhiteBoxSettings)
	require.NotNil(t, got.WhiteBoxSettings.SearchForRaceConditionsEnabled)
	require.True(t, *got.WhiteBoxSettings.SearchForRaceConditionsEnabled)
	require.NotNil(t, got.WhiteBoxSettings.AsintEnabled)
	require.True(t, *got.WhiteBoxSettings.AsintEnabled)
	require.NotNil(t, got.AsintSettings)
	require.NotNil(t, got.AsintSettings.UseNaira)
	require.True(t, *got.AsintSettings.UseNaira)
	require.NotNil(t, got.AutomatedAiTriageEnabled)
	require.True(t, *got.AutomatedAiTriageEnabled)
	require.NotNil(t, got.AutomatedAiTriageMinSeverity)
	require.Equal(t, v6_4.IssueLevelHigh, *got.AutomatedAiTriageMinSeverity)
}

func TestMapProjectSettingsToModel_OmitsEmptyTriageSeverity(t *testing.T) {
	t.Parallel()

	got := mapProjectSettingsToModel(&settings.ScanSettings{
		AutomatedAiTriageEnabled: true,
	})

	require.NotNil(t, got.AutomatedAiTriageEnabled)
	require.True(t, *got.AutomatedAiTriageEnabled)
	require.Nil(t, got.AutomatedAiTriageMinSeverity)
	require.NotNil(t, got.AsintSettings)
	require.NotNil(t, got.AsintSettings.UseNaira)
	require.False(t, *got.AsintSettings.UseNaira)
}
