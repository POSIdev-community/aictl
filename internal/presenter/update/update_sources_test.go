package update

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domainsettings "github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

func resetUpdateFlags() {
	projectIdFlag = ""
	branchIdFlag = ""
}

type fakeUpdateSourcesUC struct {
	called     int
	sourcePath string
	exclusions gitignore.Exclusions
	tempDir    string
}

func (f *fakeUpdateSourcesUC) Execute(_ context.Context, sourcePath string, exclusions gitignore.Exclusions, tempDir string) error {
	f.called++
	f.sourcePath, f.exclusions, f.tempDir = sourcePath, exclusions, tempDir
	return nil
}

type noopUpdateSettingsUC struct{}

func (noopUpdateSettingsUC) Execute(context.Context, domainsettings.ProjectSettingsPatch) error {
	return nil
}

type noopUpdateLanguagesUC struct{}

func (noopUpdateLanguagesUC) Execute(context.Context) error { return nil }

type noopUpdateSbomUC struct{}

func (noopUpdateSbomUC) Execute(context.Context, string) error { return nil }

type noopUpdateScaFeedsUC struct{}

func (noopUpdateScaFeedsUC) Execute(context.Context, string, string) error { return nil }

func TestUpdateSourcesCmd(t *testing.T) {
	t.Cleanup(resetUpdateFlags)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	branchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	srcDir := t.TempDir()
	excludeFile := filepath.Join(srcDir, "ignore.txt")
	require.NoError(t, os.WriteFile(excludeFile, []byte("*.tmp\n"), 0o644))
	tempDir := t.TempDir()

	t.Run("all_flags", func(t *testing.T) {
		resetUpdateFlags()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeUpdateSourcesUC{}
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		root := NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, uc), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}), projectCmd)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
			"-e", "*.log",
			"--exclude-from", excludeFile,
			"--temp-dir", tempDir,
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, srcDir, uc.sourcePath)
		require.Equal(t, []string{"*.log"}, uc.exclusions.Patterns)
		require.Equal(t, []string{excludeFile}, uc.exclusions.FromFiles)
		require.Equal(t, tempDir, uc.tempDir)
		require.Equal(t, projectID, cfg.ProjectId())
		require.Equal(t, branchID, cfg.BranchId())
	})

	t.Run("path_missing", func(t *testing.T) {
		resetUpdateFlags()
		cfg := cmdtest.MustCfgWithIDs(t, projectID, branchID)
		uc := &fakeUpdateSourcesUC{}
		projectCmd := NewUpdateProjectCmd(
			NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
			NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
			NewUpdateProjectLanguagesCmd(noopUpdateLanguagesUC{}),
		)
		root := NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, uc), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}), projectCmd)
		require.Error(t, cmdtest.Execute(t, root.Command, "sources", "/no/such/path"))
		require.Equal(t, 0, uc.called)
	})
}
