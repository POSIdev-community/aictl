package update

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

type fakeUpdateProjectSettingsUC struct {
	called int
	patch  domainsettings.ProjectSettingsPatch
}

func (f *fakeUpdateProjectSettingsUC) Execute(_ context.Context, patch domainsettings.ProjectSettingsPatch) error {
	f.called++
	f.patch = patch
	return nil
}

type noopUpdateSourcesUC struct{}

func (noopUpdateSourcesUC) Execute(context.Context, string, gitignore.Exclusions, string) error {
	return nil
}

func TestUpdateProjectSettingsCmd(t *testing.T) {
	t.Cleanup(resetUpdateFlags)
	agentID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	newRoot := func(uc UseCaseUpdateProjectSettings) *CmdUpdate {
		cfg := cmdtest.MustCfg(t)
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(uc),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		return NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, noopUpdateSourcesUC{}, noopUpdateLanguagesUC{}), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}, noopRollbackScaFeedsUC{}), projectCmd)
	}

	t.Run("requires_patch_flag", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateProjectSettingsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command, "project", "settings", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})

	t.Run("rejects_conflicting_preferred", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateProjectSettingsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command,
			"project", "settings", "-p", projectID.String(),
			"--preferred-agents-only", "--no-preferred-agents-only",
		))
		require.Equal(t, 0, uc.called)
	})

	t.Run("all_patch_flags", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateProjectSettingsUC{}
		cfg := cmdtest.MustCfg(t)
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(uc),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		root := NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, noopUpdateSourcesUC{}, noopUpdateLanguagesUC{}), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}, noopRollbackScaFeedsUC{}), projectCmd)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"project", "settings",
			"-p", projectID.String(),
			"--priority", "High",
			"--agents", agentID.String(),
			"--preferred-agents-only",
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, domainsettings.PriorityHigh, *uc.patch.Priority)
		require.Equal(t, []uuid.UUID{agentID}, *uc.patch.PreferredAgents)
		require.True(t, *uc.patch.PreferredAgentsOnly)
		require.Equal(t, projectID, cfg.ProjectId())
	})

	t.Run("no_preferred_agents_only", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateProjectSettingsUC{}
		require.NoError(t, cmdtest.Execute(t, newRoot(uc).Command,
			"project", "settings", "-p", projectID.String(), "--no-preferred-agents-only",
		))
		require.False(t, *uc.patch.PreferredAgentsOnly)
	})
}
