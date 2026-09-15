package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	old := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = old }()

	fn()
	require.NoError(t, w.Close())
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	return buf.String()
}

func TestReportCommandError_UserFacingOnly(t *testing.T) {
	wrapped := fmt.Errorf("create project: %w", apperror.NewNotFoundError("Project"))
	logFile := filepath.Join(t.TempDir(), "aictl.log")

	log, err := logger.NewLogger(logger.LevelQuiet, logFile)
	require.NoError(t, err)
	ctx := logger.ContextWithLogger(context.Background(), log)

	stderr := captureStderr(t, func() {
		code := reportCommandError(ctx, wrapped)
		require.Equal(t, ExitCodeAPI, code)
	})

	require.Equal(t, "Project not found\n", stderr)
	require.NotContains(t, stderr, "create project:")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	require.Contains(t, string(content), "Project not found")
	require.NotContains(t, string(content), "create project:")
}

func TestReportCommandError_DebugChainAtDebugLevel(t *testing.T) {
	wrapped := fmt.Errorf("create project: %w", apperror.NewNotFoundError("Project"))
	logFile := filepath.Join(t.TempDir(), "aictl.log")

	log, err := logger.NewLogger(logger.LevelDebug, logFile)
	require.NoError(t, err)
	ctx := logger.ContextWithLogger(context.Background(), log)

	stderr := captureStderr(t, func() {
		code := reportCommandError(ctx, wrapped)
		require.Equal(t, ExitCodeAPI, code)
	})

	// Fprintln user-facing (zap debug-core still bound to process stderr at logger creation)
	require.Contains(t, stderr, "Project not found")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	text := string(content)
	require.Contains(t, text, "Project not found")
	require.Contains(t, text, "create project:")
}

func TestReportCommandError_VerboseNoChainOnStderr(t *testing.T) {
	wrapped := fmt.Errorf("create project: %w", apperror.NewNotFoundError("Project"))
	logFile := filepath.Join(t.TempDir(), "aictl.log")

	log, err := logger.NewLogger(logger.LevelVerbose, logFile)
	require.NoError(t, err)
	ctx := logger.ContextWithLogger(context.Background(), log)

	stderr := captureStderr(t, func() {
		_ = reportCommandError(ctx, wrapped)
	})

	require.Equal(t, "Project not found\n", stderr)

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	require.Contains(t, string(content), "Project not found")
	require.NotContains(t, string(content), "create project:")
}

func TestReportCommandError_BadRequestBodyOnlyWithVerbose(t *testing.T) {
	body := `{"Message":"Invalid policy"}`
	wrapped := fmt.Errorf("set project policies: %w", apperror.NewBadRequestError(body))

	t.Run("quiet", func(t *testing.T) {
		log, err := logger.NewLogger(logger.LevelQuiet, "")
		require.NoError(t, err)
		ctx := logger.ContextWithLogger(context.Background(), log)

		stderr := captureStderr(t, func() {
			code := reportCommandError(ctx, wrapped)
			require.Equal(t, ExitCodeAPI, code)
		})

		require.Equal(t, "Bad Request error\n", stderr)
	})

	t.Run("verbose", func(t *testing.T) {
		log, err := logger.NewLogger(logger.LevelVerbose, "")
		require.NoError(t, err)
		ctx := logger.ContextWithLogger(context.Background(), log)

		stderr := captureStderr(t, func() {
			code := reportCommandError(ctx, wrapped)
			require.Equal(t, ExitCodeAPI, code)
		})

		require.Equal(t, "Bad Request error\n"+body+"\n", stderr)
	})
}
