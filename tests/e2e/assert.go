//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type pipelineMeta struct {
	ProjectID  string `json:"project_id"`
	BranchID   string `json:"branch_id"`
	ScanID     string `json:"scan_id"`
	AIEVersion string `json:"aie_version"`
}

func loadMeta(t *testing.T, path string) pipelineMeta {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "read meta.json")

	var meta pipelineMeta
	require.NoError(t, json.Unmarshal(data, &meta), "parse meta.json")

	return meta
}

func assertUUID(t *testing.T, value string) {
	t.Helper()

	_, err := uuid.Parse(strings.TrimSpace(value))
	require.NoError(t, err, "expected UUID, got %q", value)
}

func assertVersionMajorMinor(t *testing.T, version, expectedMajorMinor string) {
	t.Helper()

	version = strings.TrimSpace(version)
	require.NotEmpty(t, version, "aie version is empty")

	parts := strings.Split(version, ".")
	require.GreaterOrEqual(t, len(parts), 2, "aie version %q must contain major.minor", version)
	require.NotEmpty(t, parts[0], "aie version major in %q", version)
	require.NotEmpty(t, parts[1], "aie version minor in %q", version)

	actual := parts[0] + "." + parts[1]
	require.Equal(t, expectedMajorMinor, actual, "aie version major.minor")
}

func assertSARIF(t *testing.T, path string) {
	t.Helper()

	data, err := os.ReadFile(path)
	require.NoError(t, err, "read sarif.json")
	require.NotEmpty(t, data, "sarif.json is empty")

	var doc struct {
		Version string            `json:"version"`
		Schema  string            `json:"$schema"`
		Runs    []json.RawMessage `json:"runs"`
	}
	require.NoError(t, json.Unmarshal(data, &doc), "parse sarif.json")
	require.NotEmpty(t, doc.Version, "sarif.version")
	require.NotEmpty(t, doc.Schema, "sarif.$schema")
	require.NotEmpty(t, doc.Runs, "sarif.runs")
}

func AssertPipelineArtifacts(t *testing.T, workDir, standName string) {
	t.Helper()

	expectedVersion, err := StandVersion(standName)
	require.NoError(t, err)

	meta := loadMeta(t, filepath.Join(workDir, "meta.json"))
	assertUUID(t, meta.ProjectID)
	assertUUID(t, meta.BranchID)
	assertUUID(t, meta.ScanID)
	assertVersionMajorMinor(t, meta.AIEVersion, expectedVersion)
	assertSARIF(t, filepath.Join(workDir, "sarif.json"))
}

func AssertSbomPipelineArtifacts(t *testing.T, workDir, standName string) {
	t.Helper()

	expectedVersion, err := StandVersion(standName)
	require.NoError(t, err)

	meta := loadMeta(t, filepath.Join(workDir, "meta.json"))
	assertUUID(t, meta.ProjectID)
	assertUUID(t, meta.ScanID)
	assertVersionMajorMinor(t, meta.AIEVersion, expectedVersion)
}
