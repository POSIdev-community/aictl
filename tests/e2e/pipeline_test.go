//go:build e2e

package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasePipeline(t *testing.T) {
	configPath, err := ConfigPath()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("e2e: create tests/e2e/stands.local.yaml (make e2e-config)")
	}

	stands, err := LoadStands(configPath)
	if err != nil {
		t.Skipf("e2e: %v", err)
	}

	root, err := RepoRoot()
	require.NoError(t, err)

	aictlBin, err := ResolveAictlBin(root)
	require.NoError(t, err, "build aictl for e2e")
	t.Logf("using aictl binary: %s", aictlBin)

	e2eDir := filepath.Join(root, "tests", "e2e")
	scriptPath := filepath.Join(e2eDir, "run-pipeline.sh")
	fixturesDir := filepath.Join(e2eDir, "fixtures")

	for _, standName := range OrderedStandNames(stands) {
		stand := stands[standName]

		t.Run("AIE_"+standName, func(t *testing.T) {
			t.Parallel()

			aiprojVersion, err := stand.ResolveAiprojVersion(standName)
			require.NoError(t, err)

			workDir := t.TempDir()
			projectName := fmt.Sprintf("aictl-e2e-%s-%s", standName, uuid.NewString())

			cmd := exec.Command("bash", scriptPath, stand.URL, stand.Token, projectName)
			cmd.Dir = e2eDir
			cmd.Env = append(os.Environ(),
				"AICTL="+aictlBin,
				"FIXTURES_DIR="+fixturesDir,
				"WORK_DIR="+workDir,
				"AIPROJ_FIXTURE="+AiprojFixturePath(fixturesDir, aiprojVersion),
			)

			logWriter := &testLogWriter{t: t}
			cmd.Stdout = logWriter
			cmd.Stderr = logWriter

			runErr := cmd.Run()
			logWriter.Flush()
			require.NoError(t, runErr, "pipeline failed")

			AssertPipelineArtifacts(t, workDir, standName)
		})
	}
}

// testLogWriter streams command output into t.Log line by line.
type testLogWriter struct {
	t   *testing.T
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *testLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.t.Helper()
	w.buf.Write(p)
	for {
		line, err := w.buf.ReadString('\n')
		if err != nil {
			w.buf.WriteString(line)
			break
		}
		w.t.Log(string(bytes.TrimRight([]byte(line), "\r\n")))
	}
	return len(p), nil
}

func (w *testLogWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.t.Helper()
	if w.buf.Len() == 0 {
		return
	}
	w.t.Log(w.buf.String())
	w.buf.Reset()
}
