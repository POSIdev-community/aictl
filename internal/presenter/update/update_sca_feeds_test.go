package update

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeUpdateScaFeedsUC struct {
	called  int
	path    string
	version string
}

func (f *fakeUpdateScaFeedsUC) Execute(_ context.Context, path, version string) error {
	f.called++
	f.path = path
	f.version = version
	return nil
}

func TestUpdateScaFeedsCmd(t *testing.T) {
	t.Cleanup(resetUpdateFlags)

	zipPath := filepath.Join(t.TempDir(), "feeds.zip")
	require.NoError(t, os.WriteFile(zipPath, []byte("PK\x03\x04feeds"), 0o644))
	txtPath := filepath.Join(t.TempDir(), "feeds.txt")
	require.NoError(t, os.WriteFile(txtPath, []byte("not-zip"), 0o644))

	newRoot := func(uc UseCaseUpdateScaFeeds) *CmdUpdate {
		cfg := cmdtest.MustCfg(t)
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		return NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, noopUpdateSourcesUC{}), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(uc, noopRollbackScaFeedsUC{}), projectCmd)
	}

	t.Run("ok", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateScaFeedsUC{}
		require.NoError(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", zipPath, "--version", "47"))
		require.Equal(t, 1, uc.called)
		require.Equal(t, zipPath, uc.path)
		require.Equal(t, "47", uc.version)
	})

	t.Run("version_required", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateScaFeedsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", zipPath))
		require.Equal(t, 0, uc.called)
	})

	t.Run("path_missing", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateScaFeedsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", "/no/such.zip", "--version", "1"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("not_zip", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateScaFeedsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", txtPath, "--version", "1"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("directory_rejected", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateScaFeedsUC{}
		require.Error(t, cmdtest.Execute(t, newRoot(uc).Command, "sca-feeds", t.TempDir(), "--version", "1"))
		require.Equal(t, 0, uc.called)
	})
}
