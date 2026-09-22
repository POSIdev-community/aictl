package logger_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/pkg/logger"
)

func TestNewLogger_FileExcludesInfo(t *testing.T) {
	t.Parallel()

	logFile := filepath.Join(t.TempDir(), "aictl.log")
	log, err := logger.NewLogger(logger.LevelQuiet, logFile)
	require.NoError(t, err)
	require.True(t, log.HasFile())

	log.StdOut("stdout-result")
	log.StdErr("ops-error")
	log.Debugf("debug-chain")
	log.FileError("user-facing")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	text := string(content)

	require.NotContains(t, text, "stdout-result")
	require.Contains(t, text, "ops-error")
	require.NotContains(t, text, "debug-chain")
	require.Contains(t, text, "user-facing")
	// zap console encoder with timestamp
	require.True(t, strings.Contains(text, "T") || strings.Contains(text, "-"), "expected timestamp in file log")
}

func TestNewLogger_FileIncludesDebugAtLevelDebug(t *testing.T) {
	t.Parallel()

	logFile := filepath.Join(t.TempDir(), "aictl.log")
	log, err := logger.NewLogger(logger.LevelDebug, logFile)
	require.NoError(t, err)

	log.StdOut("stdout-result")
	log.StdErr("ops-error")
	log.Debugf("debug-chain")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	text := string(content)

	require.NotContains(t, text, "stdout-result")
	require.Contains(t, text, "ops-error")
	require.Contains(t, text, "debug-chain")
}

func TestFileError_NoOpWithoutLogPath(t *testing.T) {
	t.Parallel()

	log, err := logger.NewLogger(logger.LevelVerbose, "")
	require.NoError(t, err)
	require.False(t, log.HasFile())
	require.NotPanics(t, func() { log.FileError("msg") })
}

func TestNewLogger_DebugDoesNotDuplicateInfoOnStderr(t *testing.T) {
	// Not parallel: temporarily replaces os.Stderr.
	r, w, err := os.Pipe()
	require.NoError(t, err)

	oldStderr := os.Stderr
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = oldStderr
	})

	log, err := logger.NewLogger(logger.LevelDebug, "")
	require.NoError(t, err)

	log.StdOut("info-once")
	log.Debugf("debug-once")
	log.StdErr("error-once")

	require.NoError(t, w.Close())
	stderrBytes, err := io.ReadAll(r)
	require.NoError(t, err)
	stderr := string(stderrBytes)

	require.NotContains(t, stderr, "info-once")
	require.Contains(t, stderr, "debug-once")
	require.Contains(t, stderr, "error-once")
	require.Equal(t, 1, strings.Count(stderr, "error-once"))
	require.Equal(t, 1, strings.Count(stderr, "debug-once"))
}
