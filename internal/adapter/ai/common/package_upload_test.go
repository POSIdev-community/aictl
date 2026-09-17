package common_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

func TestComputePackageFileMeta(t *testing.T) {
	t.Parallel()

	payload := []byte("sca-feeds-payload")
	path := filepath.Join(t.TempDir(), "feeds.zip")
	require.NoError(t, os.WriteFile(path, payload, 0o644))

	sum := md5.Sum(payload)
	wantHash := hex.EncodeToString(sum[:])

	meta, err := common.ComputePackageFileMeta(path)
	require.NoError(t, err)
	require.Equal(t, "feeds.zip", meta.FileName)
	require.Equal(t, int64(len(payload)), meta.FileSize)
	require.Equal(t, wantHash, meta.Hash)
}

func TestComputePackageFileMeta_EmptyRejected(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "empty.zip")
	require.NoError(t, os.WriteFile(path, nil, 0o644))

	_, err := common.ComputePackageFileMeta(path)
	require.Error(t, err)
}

func TestPreparePackageMultipartBody(t *testing.T) {
	t.Parallel()

	payload := []byte("package-binary")
	path := filepath.Join(t.TempDir(), "AI.SCA.Feeds.1.zip")
	require.NoError(t, os.WriteFile(path, payload, 0o644))

	meta, err := common.ComputePackageFileMeta(path)
	require.NoError(t, err)

	body, contentType, err := common.PreparePackageMultipartBody(context.Background(), path, "47", meta, nil)
	require.NoError(t, err)
	defer func() { _ = body.Close() }()

	_, params, err := mime.ParseMediaType(contentType)
	require.NoError(t, err)
	boundary := params["boundary"]
	require.NotEmpty(t, boundary)

	reader := multipart.NewReader(body, boundary)

	part, err := reader.NextPart()
	require.NoError(t, err)
	require.Equal(t, "package", part.FormName())
	require.Equal(t, "AI.SCA.Feeds.1.zip", part.FileName())
	gotPackage, err := io.ReadAll(part)
	require.NoError(t, err)
	require.Equal(t, payload, gotPackage)
	_ = part.Close()

	fields := map[string]string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		value, err := io.ReadAll(part)
		require.NoError(t, err)
		fields[part.FormName()] = string(value)
		_ = part.Close()
	}

	require.Equal(t, "47", fields["version"])
	require.Equal(t, "AI.SCA.Feeds.1.zip", fields["fileName"])
	require.Equal(t, strconv.FormatInt(meta.FileSize, 10), fields["fileSize"])
	require.Equal(t, meta.Hash, fields["hash"])
}

type mockProgress struct {
	mock.Mock
	mu sync.Mutex
}

func (m *mockProgress) Report(percent int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Called(percent)
}

func TestPreparePackageMultipartBody_ReportsProgress(t *testing.T) {
	t.Parallel()

	// 10 MiB → 1 MiB read chunks map to 10% steps.
	payload := bytes.Repeat([]byte("x"), 10*1024*1024)
	path := filepath.Join(t.TempDir(), "feeds.zip")
	require.NoError(t, os.WriteFile(path, payload, 0o644))

	meta, err := common.ComputePackageFileMeta(path)
	require.NoError(t, err)

	progress := &mockProgress{}
	progress.On("Report", mock.AnythingOfType("int")).Return()

	body, _, err := common.PreparePackageMultipartBody(
		context.Background(),
		path,
		"47",
		meta,
		progress.Report,
	)
	require.NoError(t, err)
	defer func() { _ = body.Close() }()

	_, err = io.Copy(io.Discard, body)
	require.NoError(t, err)

	progress.mu.Lock()
	defer progress.mu.Unlock()
	progress.AssertExpectations(t)

	var percents []int
	for _, call := range progress.Calls {
		percents = append(percents, call.Arguments.Get(0).(int))
	}
	require.GreaterOrEqual(t, len(percents), 3)
	require.Equal(t, 0, percents[0])
	require.Equal(t, 100, percents[len(percents)-1])

	hasMid := false
	for _, p := range percents {
		if p > 0 && p < 100 {
			hasMid = true
			break
		}
	}
	require.True(t, hasMid, "expected mid-upload progress, got %v", percents)
	require.Contains(t, percents, 10)
	require.Contains(t, percents, 50)
	require.Contains(t, percents, 90)
}

func TestNewProgressReporter_StepsOfTen(t *testing.T) {
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

	report := common.NewProgressReporter(log, true, "uploading sca feeds")
	report(0)
	report(5)
	report(9)
	report(10)
	report(19)
	report(20)
	report(99)
	report(100)
	report(100)

	disabled := common.NewProgressReporter(log, false, "uploading sca feeds")
	disabled(50)

	out := buf.String()
	require.Equal(t, 1, strings.Count(out, "uploading sca feeds: 0%"))
	require.Equal(t, 1, strings.Count(out, "uploading sca feeds: 10%"))
	require.Equal(t, 1, strings.Count(out, "uploading sca feeds: 20%"))
	require.Equal(t, 1, strings.Count(out, "uploading sca feeds: 90%"))
	require.Equal(t, 1, strings.Count(out, "uploading sca feeds: 100%"))
	require.NotContains(t, out, "uploading sca feeds: 5%")
	require.NotContains(t, out, "uploading sca feeds: 50%")
}

func TestPreparePackageMultipartBody_ProgressDecadesViaMock(t *testing.T) {
	t.Parallel()

	payload := bytes.Repeat([]byte("y"), 10*1024*1024)
	path := filepath.Join(t.TempDir(), "feeds.zip")
	require.NoError(t, os.WriteFile(path, payload, 0o644))

	meta, err := common.ComputePackageFileMeta(path)
	require.NoError(t, err)

	progress := &mockProgress{}
	for _, step := range []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100} {
		progress.On("Report", step).Once()
	}

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
	reporter := common.NewProgressReporter(log, true, "uploading sca feeds")

	var (
		seenMu sync.Mutex
		seen   = map[int]struct{}{}
	)
	onProgress := func(percent int) {
		reporter(percent)
		step := percent / 10 * 10
		seenMu.Lock()
		_, already := seen[step]
		if !already {
			seen[step] = struct{}{}
		}
		seenMu.Unlock()
		if !already {
			progress.Report(step)
		}
	}

	body, _, err := common.PreparePackageMultipartBody(context.Background(), path, "1", meta, onProgress)
	require.NoError(t, err)
	defer func() { _ = body.Close() }()

	_, err = io.Copy(io.Discard, body)
	require.NoError(t, err)

	progress.AssertExpectations(t)

	out := buf.String()
	for _, step := range []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100} {
		require.Contains(t, out, "uploading sca feeds: "+strconv.Itoa(step)+"%")
	}
}
