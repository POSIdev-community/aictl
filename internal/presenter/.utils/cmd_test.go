package _utils

import (
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func resetConnectionFlags() {
	uri = ""
	token = ""
	tlsSkip = false
	cacert = ""
	verboseFlag = false
	debugFlag = false
	logPath = ""
}

func TestUpdateConnectionConfig(t *testing.T) {
	t.Cleanup(resetConnectionFlags)

	t.Run("overlays_non_empty", func(t *testing.T) {
		resetConnectionFlags()
		uri = "https://ai.example"
		token = "tok"
		tlsSkip = true

		cfg := config.NewConfig(config.Uri{}, "old", false, uuid.Nil, uuid.Nil)
		require.NoError(t, UpdateConnectionConfig(cfg))
		require.Equal(t, "https://ai.example", cfg.UriString())
		require.Equal(t, "tok", cfg.Token())
		require.True(t, cfg.TLSSkip())
	})

	t.Run("empty_flags_keep_config", func(t *testing.T) {
		resetConnectionFlags()

		uriVal, err := config.NewUri("https://keep.example")
		require.NoError(t, err)
		cfg := config.NewConfig(uriVal, "keep-token", false, uuid.Nil, uuid.Nil)

		require.NoError(t, UpdateConnectionConfig(cfg))
		require.Equal(t, "https://keep.example", cfg.UriString())
		require.Equal(t, "keep-token", cfg.Token())
		require.False(t, cfg.TLSSkip())
	})

	t.Run("overlays_cacert", func(t *testing.T) {
		resetConnectionFlags()
		cacert = "/tmp/ca.pem"

		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, UpdateConnectionConfig(cfg))
		require.Equal(t, "/tmp/ca.pem", cfg.CACertPath())
	})

	t.Run("rejects_cacert_with_tls_skip", func(t *testing.T) {
		resetConnectionFlags()
		tlsSkip = true
		cacert = "/tmp/ca.pem"

		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.Error(t, UpdateConnectionConfig(cfg))
	})
}

func TestUpdateConfig(t *testing.T) {
	t.Cleanup(resetConnectionFlags)

	t.Run("validate_fails_without_uri_token", func(t *testing.T) {
		resetConnectionFlags()
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		err := UpdateConfig(cfg)(&cobra.Command{}, nil)
		require.Error(t, err)
	})

	t.Run("ok_after_overlay", func(t *testing.T) {
		resetConnectionFlags()
		uri = "https://ai.example"
		token = "tok"
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, UpdateConfig(cfg)(&cobra.Command{}, nil))
	})
}

func TestChainRunE(t *testing.T) {
	t.Run("runs_all", func(t *testing.T) {
		var order []int
		run := ChainRunE(
			func(*cobra.Command, []string) error {
				order = append(order, 1)
				return nil
			},
			func(*cobra.Command, []string) error {
				order = append(order, 2)
				return nil
			},
		)
		require.NoError(t, run(&cobra.Command{}, nil))
		require.Equal(t, []int{1, 2}, order)
	})

	t.Run("stops_on_first_error", func(t *testing.T) {
		want := errors.New("boom")
		var secondCalled bool
		run := ChainRunE(
			func(*cobra.Command, []string) error { return want },
			func(*cobra.Command, []string) error {
				secondCalled = true
				return nil
			},
		)
		err := run(&cobra.Command{}, nil)
		require.ErrorIs(t, err, want)
		require.False(t, secondCalled)
	})
}

func TestConnectionFlagsViaCobra(t *testing.T) {
	t.Cleanup(resetConnectionFlags)
	resetConnectionFlags()

	logFile := filepath.Join(t.TempDir(), "aictl.log")
	uriVal, err := config.NewUri("https://keep.example")
	require.NoError(t, err)
	cfg := config.NewConfig(uriVal, "keep", false, uuid.Nil, uuid.Nil)

	cmd := &cobra.Command{
		Use:  "test",
		RunE: UpdateConfig(cfg),
	}
	AddConnectionPersistentFlags(cmd)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{
		"-u", "https://ai.example",
		"-t", "tok",
		"--tls-skip",
		"-v",
		"-l", logFile,
	})

	require.NoError(t, cmd.Execute())
	require.Equal(t, "https://ai.example", cfg.UriString())
	require.Equal(t, "tok", cfg.Token())
	require.True(t, cfg.TLSSkip())
	require.True(t, verboseFlag)
	require.Equal(t, logFile, logPath)
}

func TestDebugFlagViaCobra(t *testing.T) {
	t.Cleanup(resetConnectionFlags)
	resetConnectionFlags()

	cmd := &cobra.Command{Use: "test", Run: func(*cobra.Command, []string) {}}
	AddConnectionPersistentFlags(cmd)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"-V"})

	require.NoError(t, cmd.Execute())
	require.True(t, debugFlag)
	require.False(t, verboseFlag)
	require.Equal(t, 2, VerboseLevel())
}

func TestVerboseLevel_MaxOfVerboseAndDebug(t *testing.T) {
	t.Cleanup(resetConnectionFlags)

	resetConnectionFlags()
	require.Equal(t, 0, VerboseLevel())

	verboseFlag = true
	require.Equal(t, 1, VerboseLevel())

	debugFlag = true
	require.Equal(t, 2, VerboseLevel())

	verboseFlag = false
	require.Equal(t, 2, VerboseLevel())
}

func TestDebugLongFlagViaCobra(t *testing.T) {
	t.Cleanup(resetConnectionFlags)
	resetConnectionFlags()

	cmd := &cobra.Command{Use: "test", Run: func(*cobra.Command, []string) {}}
	AddConnectionPersistentFlags(cmd)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"--debug", "--verbose"})

	require.NoError(t, cmd.Execute())
	require.True(t, debugFlag)
	require.True(t, verboseFlag)
	require.Equal(t, 2, VerboseLevel())
}
