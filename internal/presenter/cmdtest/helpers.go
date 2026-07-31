package cmdtest

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func MustCfg(t *testing.T) *config.Config {
	t.Helper()
	uri, err := config.NewUri("https://ai.example")
	require.NoError(t, err)
	return config.NewConfig(uri, "test-token", true, uuid.Nil, uuid.Nil)
}

func MustCfgWithIDs(t *testing.T, projectID, branchID uuid.UUID) *config.Config {
	t.Helper()
	uri, err := config.NewUri("https://ai.example")
	require.NoError(t, err)
	return config.NewConfig(uri, "test-token", true, projectID, branchID)
}

func SilenceCmd(cmd *cobra.Command) {
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
}

func WithStdin(t *testing.T, input string, fn func()) {
	t.Helper()
	old := os.Stdin
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = old
		_ = r.Close()
	})
	_, err = io.Copy(w, strings.NewReader(input))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	fn()
}

func WithStdio(t *testing.T, fn func()) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	require.NoError(t, err)
	os.Stdout = devNull
	os.Stderr = devNull
	t.Cleanup(func() {
		os.Stdout = oldOut
		os.Stderr = oldErr
		_ = devNull.Close()
	})
	fn()
}

func Execute(t *testing.T, cmd *cobra.Command, args ...string) error {
	t.Helper()
	SilenceCmd(cmd)
	cmd.SetArgs(args)
	var execErr error
	WithStdio(t, func() {
		execErr = cmd.Execute()
	})
	return execErr
}
