package create

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

func resetSafeFlag() { safeFlag = false }

type fakeCreateBranchUC struct {
	called     int
	branchName string
	scanTarget string
	safe       bool
	exclusions gitignore.Exclusions
	tempDir    string
}

func (f *fakeCreateBranchUC) Execute(
	_ context.Context,
	_ *config.Config,
	branchName, scanTarget string,
	safe bool,
	exclusions gitignore.Exclusions,
	tempDir string,
) error {
	f.called++
	f.branchName, f.scanTarget, f.safe = branchName, scanTarget, safe
	f.exclusions, f.tempDir = exclusions, tempDir
	return nil
}

type noopCreateProjectUC struct{}

func (noopCreateProjectUC) Execute(context.Context, string, bool) error { return nil }

type noopCreateSbomProjectUC struct{}

func (noopCreateSbomProjectUC) Execute(context.Context, string, string, bool) error { return nil }

func TestCreateBranchCmd(t *testing.T) {
	t.Cleanup(resetSafeFlag)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	srcDir := t.TempDir()
	excludeFile := filepath.Join(srcDir, "ignore.txt")
	require.NoError(t, os.WriteFile(excludeFile, []byte("*.tmp\n"), 0o644))
	tempDir := t.TempDir()

	t.Run("all_flags", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateBranchUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, uc), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(noopCreateSbomProjectUC{}))
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"branch", "feature-x",
			"-p", projectID.String(),
			"-s", srcDir,
			"-e", "*.log",
			"--exclude-from", excludeFile,
			"--temp-dir", tempDir,
			"--safe",
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, "feature-x", uc.branchName)
		require.Equal(t, srcDir, uc.scanTarget)
		require.True(t, uc.safe)
		require.Equal(t, []string{"*.log"}, uc.exclusions.Patterns)
		require.Equal(t, []string{excludeFile}, uc.exclusions.FromFiles)
		require.Equal(t, tempDir, uc.tempDir)
		require.Equal(t, projectID, cfg.ProjectId())
	})

	t.Run("invalid_scan_target", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateBranchUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, uc), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(noopCreateSbomProjectUC{}))
		require.Error(t, cmdtest.Execute(t, root.Command, "branch", "main", "-s", "/no/such/path", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})

	t.Run("stdin_dash", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateBranchUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, uc), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(noopCreateSbomProjectUC{}))
		cmdtest.WithStdin(t, "branch-from-stdin\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "branch", "-", "-p", projectID.String()))
		})
		require.Equal(t, "branch-from-stdin", uc.branchName)
	})

	t.Run("missing_name", func(t *testing.T) {
		resetSafeFlag()
		cfg := cmdtest.MustCfg(t)
		uc := &fakeCreateBranchUC{}
		root := NewCreateCmd(cfg, NewCreateBranchCmd(cfg, uc), NewCreateProjectCmd(noopCreateProjectUC{}), NewCreateSbomProjectCmd(noopCreateSbomProjectUC{}))
		err := cmdtest.Execute(t, root.Command, "branch", "-p", projectID.String())
		require.Error(t, err)
		var requiredErr *validation.RequiredError
		require.True(t, errors.As(err, &requiredErr))
		require.Equal(t, "branch-name", requiredErr.Field)
		require.Equal(t, 0, uc.called)
	})
}
