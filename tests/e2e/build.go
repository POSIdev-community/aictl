//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveAictlBin returns AICTL_BIN if set; otherwise builds aictl into bin/aictl.
func ResolveAictlBin(root string) (string, error) {
	if bin := os.Getenv("AICTL_BIN"); bin != "" {
		if _, err := os.Stat(bin); err != nil {
			return "", fmt.Errorf("AICTL_BIN %q: %w", bin, err)
		}

		return bin, nil
	}

	return BuildAictl(root)
}

// BuildAictl compiles cmd/run into bin/aictl (same as make build-e2e).
func BuildAictl(root string) (string, error) {
	out := filepath.Join(root, "bin", "aictl")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return "", fmt.Errorf("mkdir bin: %w", err)
	}

	version := "dev"
	if data, err := os.ReadFile(filepath.Join(root, "VERSION")); err == nil {
		version = strings.TrimSpace(string(data))
	}

	ldflags := fmt.Sprintf("-X 'github.com/POSIdev-community/aictl/pkg/version.version=%s' -s -w", version)
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-trimpath", "-o", out, "cmd/run/main.go")
	cmd.Dir = root
	if outBytes, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build aictl: %w\n%s", err, outBytes)
	}

	return out, nil
}
