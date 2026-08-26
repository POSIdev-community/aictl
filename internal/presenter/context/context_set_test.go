package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeConfigSetUC struct{ called int }

func (f *fakeConfigSetUC) Execute() error { f.called++; return nil }

func TestConfigSetCommand(t *testing.T) {
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	branchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	t.Run("requires_at_least_one_flag", func(t *testing.T) {
		cfg := emptyCfg()
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.Error(t, cmdtest.Execute(t, root.Command, "set"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("rejects_both_tls_flags", func(t *testing.T) {
		cfg := emptyCfg()
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.Error(t, cmdtest.Execute(t, root.Command, "set", "--tls-skip", "--no-tls-skip"))
		require.Equal(t, 0, uc.called)
	})

	t.Run("overlays_all_flags_and_executes", func(t *testing.T) {
		cfg := emptyCfg()
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.NoError(t, cmdtest.Execute(t, root.Command,
			"set",
			"-u", "https://ai.example",
			"-t", "token",
			"--tls-skip",
			"-p", projectID.String(),
			"-b", branchID.String(),
		))
		require.Equal(t, 1, uc.called)
		require.Equal(t, "https://ai.example", cfg.UriString())
		require.Equal(t, "token", cfg.Token())
		require.True(t, cfg.TLSSkip())
		require.Equal(t, projectID, cfg.ProjectId())
		require.Equal(t, branchID, cfg.BranchId())
	})

	t.Run("no_tls_skip", func(t *testing.T) {
		cfg := emptyCfg()
		cfg.SetTLSSkip(true)
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.NoError(t, cmdtest.Execute(t, root.Command, "set", "--no-tls-skip"))
		require.Equal(t, 1, uc.called)
		require.False(t, cfg.TLSSkip())
	})

	t.Run("sets_cacert", func(t *testing.T) {
		ca := filepath.Join(t.TempDir(), "ca.pem")
		require.NoError(t, os.WriteFile(ca, []byte("-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----\n"), 0o600))

		cfg := emptyCfg()
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.NoError(t, cmdtest.Execute(t, root.Command, "set", "--cacert", ca))
		require.Equal(t, 1, uc.called)
		require.Equal(t, ca, cfg.CACertPath())
	})

	t.Run("rejects_cacert_with_tls_skip", func(t *testing.T) {
		ca := filepath.Join(t.TempDir(), "ca.pem")
		require.NoError(t, os.WriteFile(ca, []byte("x"), 0o600))

		cfg := emptyCfg()
		uc := &fakeConfigSetUC{}
		root := NewContextCmd(NewConfigClearCommand(noopClearUC{}), NewConfigSetCommand(cfg, uc), NewConfigShowCommand(noopShowUC{}), NewConfigUnsetCommand(noopUnsetUC{}))
		require.Error(t, cmdtest.Execute(t, root.Command, "set", "--cacert", ca, "--tls-skip"))
		require.Equal(t, 0, uc.called)
	})
}
