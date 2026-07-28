package gitignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/pkg/gitignore"
)

func TestMatcher_multiplePatterns(t *testing.T) {
	t.Parallel()

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		Patterns: []string{".git", "node_modules/", "*.log"},
	})
	require.NoError(t, err)

	assert.True(t, matcher.Match(".git", true))
	assert.True(t, matcher.Match("node_modules/pkg/index.js", false))
	assert.True(t, matcher.Match("debug.log", false))
	assert.False(t, matcher.Match("src/main.go", false))
}

func TestMatcher_negation(t *testing.T) {
	t.Parallel()

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		Patterns: []string{"dir*", "!dirname"},
	})
	require.NoError(t, err)

	assert.False(t, matcher.Match("dirname", true))
	assert.False(t, matcher.Match("dirname/foo", false))
	assert.True(t, matcher.Match("dirfoo", true))
}

func TestMatcher_directoryOnlyPattern(t *testing.T) {
	t.Parallel()

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		Patterns: []string{"dirname/"},
	})
	require.NoError(t, err)

	assert.True(t, matcher.Match("dirname", true))
	assert.False(t, matcher.Match("dirname", false))
}

func TestMatcher_multipleFromFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	first := filepath.Join(dir, "first.ignore")
	second := filepath.Join(dir, "second.ignore")
	require.NoError(t, os.WriteFile(first, []byte("*.tmp\n"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("build/\n"), 0o644))

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		FromFiles: []string{first, second},
	})
	require.NoError(t, err)

	assert.True(t, matcher.Match("file.tmp", false))
	assert.True(t, matcher.Match("build", true))
	assert.False(t, matcher.Match("src/main.go", false))
}

func TestMatcher_patternsAndFilesCombined(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "excludes.txt")
	require.NoError(t, os.WriteFile(path, []byte("*.tmp\n"), 0o644))

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		Patterns:  []string{".git"},
		FromFiles: []string{path},
	})
	require.NoError(t, err)

	assert.True(t, matcher.Match(".git", true))
	assert.True(t, matcher.Match("cache.tmp", false))
	assert.False(t, matcher.Match("src/main.go", false))
}

func TestCollectReader_orderPreservesNegation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "excludes.txt")
	require.NoError(t, os.WriteFile(path, []byte("!important.log\n"), 0o644))

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{
		Patterns:  []string{"*.log"},
		FromFiles: []string{path},
	})
	require.NoError(t, err)

	assert.True(t, matcher.Match("debug.log", false))
	assert.False(t, matcher.Match("important.log", false))
}

func TestCompile_EmptyExclusions(t *testing.T) {
	t.Parallel()

	matcher, err := gitignore.NewMatcher(gitignore.Exclusions{})
	require.NoError(t, err)
	assert.False(t, matcher.Match(".git", true))
}
