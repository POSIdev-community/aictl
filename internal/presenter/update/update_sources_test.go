package update

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
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
	err        error
	onCall     func()
}

func (f *fakeUpdateSourcesUC) Execute(_ context.Context, sourcePath string, exclusions gitignore.Exclusions, tempDir string) error {
	f.called++
	f.sourcePath, f.exclusions, f.tempDir = sourcePath, exclusions, tempDir
	if f.onCall != nil {
		f.onCall()
	}
	return f.err
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

type noopRollbackScaFeedsUC struct{}

func (noopRollbackScaFeedsUC) Execute(context.Context, bool) error { return nil }

func newUpdateSourcesRoot(t *testing.T, cfg *config.Config, sourcesUC UseCaseUpdateSources, languagesUC UseCaseUpdateProjectLanguages) *CmdUpdate {
	t.Helper()
	projectCmd := NewUpdateProjectCmd(
		NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
		NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
		NewUpdateProjectLanguagesCmd(languagesUC),
	)
	return NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, sourcesUC, languagesUC), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}, noopRollbackScaFeedsUC{}), projectCmd)
}

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
		lang := &fakeUpdateLanguagesUC{}
		root := newUpdateSourcesRoot(t, cfg, uc, lang)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
			"-e", "*.log",
			"--exclude-from", excludeFile,
			"--temp-dir", tempDir,
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, 0, lang.called)
		require.Equal(t, srcDir, uc.sourcePath)
		require.Equal(t, []string{"*.log"}, uc.exclusions.Patterns)
		require.Equal(t, []string{excludeFile}, uc.exclusions.FromFiles)
		require.Equal(t, tempDir, uc.tempDir)
		require.Equal(t, projectID, cfg.ProjectId())
		require.Equal(t, branchID, cfg.BranchId())
	})

	t.Run("update_languages_off_by_default", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateSourcesUC{}
		lang := &fakeUpdateLanguagesUC{}
		root := newUpdateSourcesRoot(t, cmdtest.MustCfg(t), uc, lang)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, 0, lang.called)
	})

	t.Run("update_languages_on_calls_languages_after_sources", func(t *testing.T) {
		resetUpdateFlags()
		var order []string
		uc := &fakeUpdateSourcesUC{onCall: func() { order = append(order, "sources") }}
		lang := &fakeUpdateLanguagesUC{onCall: func() { order = append(order, "languages") }}
		root := newUpdateSourcesRoot(t, cmdtest.MustCfg(t), uc, lang)
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
			"--update-languages",
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, 1, lang.called)
		require.Equal(t, []string{"sources", "languages"}, order)
	})

	t.Run("sources_error_skips_languages", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateSourcesUC{err: errors.New("upload failed")}
		lang := &fakeUpdateLanguagesUC{}
		root := newUpdateSourcesRoot(t, cmdtest.MustCfg(t), uc, lang)
		err := cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
			"--update-languages",
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "update sources")
		require.Equal(t, 1, uc.called)
		require.Equal(t, 0, lang.called)
	})

	t.Run("languages_error_propagates", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateSourcesUC{}
		lang := &fakeUpdateLanguagesUC{err: errors.New("languages failed")}
		root := newUpdateSourcesRoot(t, cmdtest.MustCfg(t), uc, lang)
		err := cmdtest.Execute(t, root.Command,
			"sources", srcDir,
			"-p", projectID.String(),
			"-b", branchID.String(),
			"--update-languages",
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "update project languages")
		require.Equal(t, 1, uc.called)
		require.Equal(t, 1, lang.called)
	})

	t.Run("path_missing", func(t *testing.T) {
		resetUpdateFlags()
		uc := &fakeUpdateSourcesUC{}
		lang := &fakeUpdateLanguagesUC{}
		root := newUpdateSourcesRoot(t, cmdtest.MustCfgWithIDs(t, projectID, branchID), uc, lang)
		require.Error(t, cmdtest.Execute(t, root.Command, "sources", "/no/such/path"))
		require.Equal(t, 0, uc.called)
		require.Equal(t, 0, lang.called)
	})
}
