package download_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scafeeds/download"
)

type failAfterReader struct {
	data []byte
	read int
	err  error
}

func (r *failAfterReader) Read(p []byte) (int, error) {
	if r.read >= len(r.data) {
		return 0, r.err
	}
	n := copy(p, r.data[r.read:])
	r.read += n
	if r.read >= len(r.data) {
		return n, r.err
	}

	return n, nil
}

func (r *failAfterReader) Close() error { return nil }

func TestDownloadScaFeedsUseCase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("write_to_named_file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		dest := filepath.Join(dir, "out.zip")

		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.2.3").Return(io.NopCloser(bytes.NewReader([]byte("zip"))), "feeds.zip", nil).Once()
		cli.On("ShowTextf", ctx, "saved sca feeds to '%s'", mock.Anything).Return().Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, "1.2.3", dest))

		data, err := os.ReadFile(dest)
		require.NoError(t, err)
		require.Equal(t, []byte("zip"), data)
	})

	t.Run("write_into_dir", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.0").Return(io.NopCloser(bytes.NewReader([]byte("data"))), "a.zip", nil).Once()
		cli.On("ShowTextf", ctx, "saved sca feeds to '%s'", mock.Anything).Return().Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, "1.0", dir))

		data, err := os.ReadFile(filepath.Join(dir, "a.zip"))
		require.NoError(t, err)
		require.Equal(t, []byte("data"), data)
	})

	t.Run("default_filename_when_empty", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "9.9").Return(io.NopCloser(bytes.NewReader([]byte("z"))), "", nil).Once()
		cli.On("ShowTextf", ctx, "saved sca feeds to '%s'", mock.Anything).Return().Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, "9.9", dir))

		data, err := os.ReadFile(filepath.Join(dir, "sca-feeds-9.9.zip"))
		require.NoError(t, err)
		require.Equal(t, []byte("z"), data)
	})

	t.Run("existing_file_error", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "out.zip")
		require.NoError(t, os.WriteFile(path, []byte("x"), 0o644))

		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.0").Return(io.NopCloser(bytes.NewReader([]byte("data"))), "a.zip", nil).Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		err = uc.Execute(ctx, "1.0", path)
		require.Error(t, err)
		var ve *validation.Error
		require.ErrorAs(t, err, &ve)
	})

	t.Run("existing_file_in_dir_error", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "a.zip"), []byte("x"), 0o644))

		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.0").Return(io.NopCloser(bytes.NewReader([]byte("data"))), "a.zip", nil).Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.Error(t, uc.Execute(ctx, "1.0", dir))
	})

	t.Run("interrupted_download_removes_partial", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		dest := filepath.Join(dir, "feeds.zip")

		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.2.3").Return(
			&failAfterReader{data: []byte("partial"), err: errors.New("interrupted")},
			"feeds.zip",
			nil,
		).Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.Error(t, uc.Execute(ctx, "1.2.3", dest))

		_, statErr := os.Stat(dest)
		require.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("canceled_download_removes_partial_in_dir", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()

		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "3.0").Return(
			&failAfterReader{data: bytes.Repeat([]byte("p"), 128), err: context.Canceled},
			"feeds.zip",
			nil,
		).Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		err = uc.Execute(ctx, "3.0", dir)
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)

		_, statErr := os.Stat(filepath.Join(dir, "feeds.zip"))
		require.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("download_api_error", func(t *testing.T) {
		t.Parallel()
		ai := download.NewMockAI(t)
		cli := download.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("DownloadScaFeeds", ctx, "1.0").Return(nil, "", errors.New("boom")).Once()

		uc, err := download.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.Error(t, uc.Execute(ctx, "1.0", filepath.Join(t.TempDir(), "x.zip")))
	})

	t.Run("nil_adapters", func(t *testing.T) {
		t.Parallel()
		_, err := download.NewUseCase(nil, download.NewMockCLI(t))
		require.Error(t, err)
		_, err = download.NewUseCase(download.NewMockAI(t), nil)
		require.Error(t, err)
	})
}
