//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// RunAictl runs aictl with stand connection flags (-u/-t/--tls-skip) and requires exit 0.
// Connection flags are inserted after the first argument (create|get|scan|set|update|delete),
// matching run-pipeline.sh.
func RunAictl(t *testing.T, bin string, stand Stand, env []string, args ...string) string {
	t.Helper()

	stdout, err := runAictl(t, bin, stand, env, args...)
	require.NoError(t, err, "aictl %s", strings.Join(args, " "))

	return stdout
}

// RunAictlNoConn runs aictl without injecting connection flags (for ctx commands).
func RunAictlNoConn(t *testing.T, bin string, env []string, args ...string) string {
	t.Helper()

	stdout, err := runAictlArgs(t, bin, env, args...)
	require.NoError(t, err, "aictl %s", strings.Join(args, " "))

	return stdout
}

func runAictl(t *testing.T, bin string, stand Stand, env []string, args ...string) (string, error) {
	t.Helper()

	require.NotEmpty(t, args, "aictl args")

	conn := []string{"-u", stand.URL, "-t", stand.Token, "--tls-skip"}
	fullArgs := make([]string, 0, len(args)+len(conn))
	fullArgs = append(fullArgs, args[0])
	fullArgs = append(fullArgs, conn...)
	fullArgs = append(fullArgs, args[1:]...)

	return runAictlArgs(t, bin, env, fullArgs...)
}

func runAictlArgs(t *testing.T, bin string, env []string, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)

	var stdout bytes.Buffer
	logWriter := &testLogWriter{t: t}
	cmd.Stdout = &stdoutTee{buf: &stdout, log: logWriter}
	cmd.Stderr = logWriter

	err := cmd.Run()
	logWriter.Flush()

	return strings.TrimSpace(stdout.String()), err
}

type stdoutTee struct {
	buf *bytes.Buffer
	log *testLogWriter
}

func (t *stdoutTee) Write(p []byte) (int, error) {
	n, err := t.buf.Write(p)
	if err != nil {
		return n, err
	}
	_, _ = t.log.Write(p)

	return n, nil
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
