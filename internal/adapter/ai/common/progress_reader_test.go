package common

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/POSIdev-community/aictl/pkg/logger"
)

func TestProgressReadCloser(t *testing.T) {
	t.Parallel()

	var percents []int
	onProgress := func(p int) { percents = append(percents, p) }

	payload := bytes.Repeat([]byte("x"), 1000)
	rc := ProgressReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), onProgress)

	buf := make([]byte, 100)
	for {
		_, err := rc.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}
	require.NoError(t, rc.Close())

	require.Equal(t, 0, percents[0])
	require.Equal(t, 100, percents[len(percents)-1])
	require.Contains(t, percents, 10)
	require.Contains(t, percents, 50)
	require.Contains(t, percents, 90)
}

func TestProgressReadCloser_NilCallback(t *testing.T) {
	t.Parallel()

	payload := []byte("data")
	rc := ProgressReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), nil)
	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	require.Equal(t, payload, data)
	require.NoError(t, rc.Close())
}

func TestProgressReadCloser_UnknownSize(t *testing.T) {
	t.Parallel()

	var percents []int
	payload := []byte("abcdef")
	rc := ProgressReadCloser(io.NopCloser(bytes.NewReader(payload)), 0, func(p int) {
		percents = append(percents, p)
	})

	_, err := io.Copy(io.Discard, rc)
	require.NoError(t, err)
	require.NoError(t, rc.Close())

	require.Equal(t, []int{0, 100}, percents)
}

func TestProgressReadCloser_WithReporter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		}),
		zapcore.AddSync(&buf),
		zapcore.ErrorLevel,
	)
	log := logger.Wrap(zap.New(core))
	report := NewProgressReporter(log, true, "downloading sca feeds")

	payload := bytes.Repeat([]byte("y"), 10*1024)
	rc := ProgressReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), report)
	readBuf := make([]byte, 256)
	for {
		_, err := rc.Read(readBuf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}
	require.NoError(t, rc.Close())

	out := buf.String()
	require.Contains(t, out, "downloading sca feeds: 0%")
	require.Contains(t, out, "downloading sca feeds: 10%")
	require.Contains(t, out, "downloading sca feeds: 50%")
	require.Contains(t, out, "downloading sca feeds: 90%")
	require.Contains(t, out, "downloading sca feeds: 100%")
	require.Equal(t, 1, strings.Count(out, "downloading sca feeds: 0%"))
	require.Equal(t, 1, strings.Count(out, "downloading sca feeds: 100%"))
	require.NotContains(t, out, "downloading sca feeds: 5%")
}

func TestProgressReadCloser_QuietLoggerNoProgress(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	// Quiet: only Info on stdout — StdErr (Error) is dropped.
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		}),
		zapcore.AddSync(&buf),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		}),
	)
	log := logger.Wrap(zap.New(core))
	report := NewProgressReporter(log, true, "downloading sca feeds")

	payload := bytes.Repeat([]byte("z"), 2048)
	rc := ProgressReadCloser(io.NopCloser(bytes.NewReader(payload)), int64(len(payload)), report)
	_, err := io.Copy(io.Discard, rc)
	require.NoError(t, err)
	require.NoError(t, rc.Close())

	require.NotContains(t, buf.String(), "downloading sca feeds")
}
