package utils_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	utils "github.com/POSIdev-community/aictl/internal/core/usecase/.utils"
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

func TestCopyFileToPath_SuccessKeepsFile(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "ok.zip")
	require.NoError(t, utils.CopyFileToPath(io.NopCloser(bytes.NewReader([]byte("feeds"))), dest))

	data, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, []byte("feeds"), data)
}

func TestCopyFileToPath_RemovesPartialOnError(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "partial.zip")
	src := &failAfterReader{data: []byte("partial"), err: errors.New("interrupted")}

	err := utils.CopyFileToPath(src, dest)
	require.Error(t, err)
	_, statErr := os.Stat(dest)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestCopyFileToPath_RemovesPartialOnCanceled(t *testing.T) {
	t.Parallel()

	dest := filepath.Join(t.TempDir(), "canceled.zip")
	src := &failAfterReader{data: bytes.Repeat([]byte("x"), 64), err: context.Canceled}

	err := utils.CopyFileToPath(src, dest)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	_, statErr := os.Stat(dest)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestCopyFileToPath_CreateError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	err := utils.CopyFileToPath(io.NopCloser(bytes.NewReader([]byte("x"))), dir)
	require.Error(t, err)
}
