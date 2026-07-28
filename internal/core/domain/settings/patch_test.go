package settings_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
)

func TestScanSettings_Patch(t *testing.T) {
	agent1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	agent2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	current := domainsettings.ScanSettings{
		ProjectName: "unchanged",
		Priority:    domainsettings.PriorityLow,
		PreferredAgentsSettings: domainsettings.PreferredAgentsSettings{
			PreferredAgents:     []uuid.UUID{agent1},
			PreferredAgentsOnly: false,
		},
	}

	priority := domainsettings.PriorityHigh
	preferredOnly := true
	agents := []uuid.UUID{agent2}

	current.Patch(domainsettings.ProjectSettingsPatch{
		Priority:            &priority,
		PreferredAgentsOnly: &preferredOnly,
		PreferredAgents:     &agents,
	})

	require.Equal(t, "unchanged", current.ProjectName)
	require.Equal(t, domainsettings.PriorityHigh, current.Priority)
	require.True(t, current.PreferredAgentsSettings.PreferredAgentsOnly)
	require.Equal(t, []uuid.UUID{agent2}, current.PreferredAgentsSettings.PreferredAgents)
}

func TestScanSettings_Patch_PreferredAgentsOnlyFalse(t *testing.T) {
	current := domainsettings.ScanSettings{
		PreferredAgentsSettings: domainsettings.PreferredAgentsSettings{
			PreferredAgentsOnly: true,
		},
	}

	preferredOnly := false
	current.Patch(domainsettings.ProjectSettingsPatch{
		PreferredAgentsOnly: &preferredOnly,
	})

	require.False(t, current.PreferredAgentsSettings.PreferredAgentsOnly)
}

func TestParsePriority(t *testing.T) {
	p, err := domainsettings.ParsePriority("High")
	require.NoError(t, err)
	require.Equal(t, domainsettings.PriorityHigh, p)

	_, err = domainsettings.ParsePriority("Invalid")
	require.Error(t, err)
}

func TestParsePreferredAgentsCSV(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	agents, err := domainsettings.ParsePreferredAgentsCSV(id.String())
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{id}, agents)

	_, err = domainsettings.ParsePreferredAgentsCSV("not-a-uuid")
	require.Error(t, err)
}
