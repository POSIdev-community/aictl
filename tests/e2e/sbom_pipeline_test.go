//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSbomPipeline(t *testing.T) {
	configPath, err := ConfigPath()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("e2e: create tests/e2e/stands.local.yaml (task e2e-config)")
	}

	stands, err := LoadStands(configPath)
	if err != nil {
		t.Skipf("e2e: %v", err)
	}

	stand, ok := stands[standOrder63]
	if !ok {
		t.Skip("e2e: stand 6.3 not configured")
	}

	root, err := RepoRoot()
	require.NoError(t, err)

	aictlBin, err := ResolveAictlBin(root)
	require.NoError(t, err, "build aictl for e2e")
	t.Logf("using aictl binary: %s", aictlBin)

	e2eDir := filepath.Join(root, "tests", "e2e")
	scriptPath := filepath.Join(e2eDir, "run-sbom-pipeline.sh")
	fixturesDir := filepath.Join(e2eDir, "fixtures")
	sbomFile := filepath.Join(fixturesDir, "sbom", "sbom.json")
	if _, err := os.Stat(sbomFile); os.IsNotExist(err) {
		t.Fatalf("SBOM fixture missing: %s — place an SBOM JSON file at tests/e2e/fixtures/sbom/sbom.json before running the SBOM pipeline", sbomFile)
	}

	workDir := t.TempDir()
	projectName := fmt.Sprintf("aictl-e2e-sbom-%s", uuid.NewString())

	cmd := exec.Command("bash", scriptPath, stand.URL, stand.Token, projectName)
	cmd.Dir = e2eDir
	cmd.Env = append(os.Environ(),
		"AICTL="+aictlBin,
		"FIXTURES_DIR="+fixturesDir,
		"WORK_DIR="+workDir,
	)

	logWriter := &testLogWriter{t: t}
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	runErr := cmd.Run()
	logWriter.Flush()
	require.NoError(t, runErr, "sbom pipeline failed")

	AssertSbomPipelineArtifacts(t, workDir, standOrder63)
}
