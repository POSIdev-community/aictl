package common_test

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

func TestPrepareArchive_excludesPatterns(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, ".git"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".git", "config"), []byte("git"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(root, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "src", "main.go"), []byte("package main"), 0o644))

	archivePath, err := common.PrepareArchive(context.Background(), root, gitignore.Exclusions{Patterns: []string{".git"}}, "")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(archivePath)
	}()

	reader, err := zip.OpenReader(archivePath)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()

	var names []string
	for _, file := range reader.File {
		names = append(names, file.Name)
	}

	assert.Contains(t, names, "src/main.go")
	assert.NotContains(t, names, ".git/config")
	assert.NotContains(t, names, ".git/")
}

func TestPrepareArchive_multipleExcludeFromFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "app.go"), []byte("package main"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "cache.tmp"), []byte("tmp"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(root, "build"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "out.bin"), []byte("bin"), 0o644))

	first := filepath.Join(root, "first.ignore")
	second := filepath.Join(root, "second.ignore")
	require.NoError(t, os.WriteFile(first, []byte("*.tmp\n"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("build/\n"), 0o644))

	archivePath, err := common.PrepareArchive(context.Background(), root, gitignore.Exclusions{FromFiles: []string{first, second}}, "")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(archivePath)
	}()

	reader, err := zip.OpenReader(archivePath)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()

	var names []string
	for _, file := range reader.File {
		names = append(names, file.Name)
	}

	assert.Contains(t, names, "app.go")
	assert.NotContains(t, names, "cache.tmp")
	assert.NotContains(t, names, "build/out.bin")
}

func TestPrepareArchive_zipSourceUnchanged(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	zipPath := filepath.Join(root, "sources.zip")
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	writer := zip.NewWriter(zipFile)
	fileWriter, err := writer.Create("main.go")
	require.NoError(t, err)
	_, err = io.WriteString(fileWriter, "package main")
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	require.NoError(t, zipFile.Close())

	archivePath, err := common.PrepareArchive(context.Background(), zipPath, gitignore.Exclusions{Patterns: []string{".git"}}, "")
	require.NoError(t, err)
	assert.Equal(t, zipPath, archivePath)
}

func TestPrepareArchive_tempDir(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0o644))

	tempDir := t.TempDir()
	archivePath, err := common.PrepareArchive(context.Background(), root, gitignore.Exclusions{}, tempDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(archivePath)
	}()

	assert.Equal(t, tempDir, filepath.Dir(archivePath))
	assert.True(t, strings.HasSuffix(archivePath, ".zip"))
}
