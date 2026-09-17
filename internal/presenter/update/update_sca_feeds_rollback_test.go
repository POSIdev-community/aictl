package update

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeRollbackScaFeedsUC struct {
	called      int
	skipConfirm bool
}

func (f *fakeRollbackScaFeedsUC) Execute(_ context.Context, skipConfirm bool) error {
	f.called++
	f.skipConfirm = skipConfirm
	return nil
}

func TestUpdateScaFeedsRollbackCmd(t *testing.T) {
	t.Cleanup(resetUpdateFlags)

	newRoot := func(rb UseCaseRollbackScaFeeds) *CmdUpdate {
		cfg := cmdtest.MustCfg(t)
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		return NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, noopUpdateSourcesUC{}), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}, rb), projectCmd)
	}

	t.Run("with_yes", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeRollbackScaFeedsUC{}
		require.NoError(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", "rollback", "-y"))
		require.Equal(t, 1, uc.called)
		require.True(t, uc.skipConfirm)
	})

	t.Run("without_yes", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeRollbackScaFeedsUC{}
		require.NoError(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", "rollback"))
		require.Equal(t, 1, uc.called)
		require.False(t, uc.skipConfirm)
	})
}
