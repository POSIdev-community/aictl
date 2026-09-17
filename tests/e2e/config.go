//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigRelPath = "tests/e2e/stands.local.yaml"
	standOrder53         = "5.3"
	standOrder54         = "5.4"
	standOrder60         = "6.0"
	standOrder61         = "6.1"
	standOrder62         = "6.2"
	standOrder63         = "6.3"
)

var standOrder = []string{standOrder53, standOrder54, standOrder60, standOrder61, standOrder62, standOrder63}

type Stand struct {
	URL           string `yaml:"url"`
	Token         string `yaml:"token"`
	AiprojVersion string `yaml:"aiproj_version"`
}

type standsFile struct {
	Stands  map[string]Stand `yaml:"stands"`
	TLSSkip bool             `yaml:"tls_skip"`
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %q", wd)
		}
		dir = parent
	}
}

func ConfigPath() (string, error) {
	if path := os.Getenv("AICTL_E2E_CONFIG"); path != "" {
		return path, nil
	}

	root, err := RepoRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, defaultConfigRelPath), nil
}

func LoadStands(path string) (map[string]Stand, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg standsFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	if len(cfg.Stands) == 0 {
		return nil, fmt.Errorf("config %q: no stands defined", path)
	}

	for name, stand := range cfg.Stands {
		if strings.TrimSpace(stand.URL) == "" {
			return nil, fmt.Errorf("stand %q: url is required", name)
		}
		if strings.TrimSpace(stand.Token) == "" || stand.Token == "<token>" {
			return nil, fmt.Errorf("stand %q: token is not configured", name)
		}
		if _, err := StandVersion(name); err != nil {
			return nil, err
		}
	}

	return cfg.Stands, nil
}

func StandVersion(standName string) (string, error) {
	parts := strings.Split(strings.TrimSpace(standName), ".")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("stand %q: name must be major.minor (e.g. 6.0)", standName)
	}

	return parts[0] + "." + parts[1], nil
}

func (s Stand) ResolveAiprojVersion(standName string) (string, error) {
	if v := strings.TrimSpace(s.AiprojVersion); v != "" {
		return v, nil
	}

	switch standName {
	case standOrder53, standOrder54:
		return "1.9", nil
	case standOrder60:
		return "1.10", nil
	case standOrder61, standOrder62:
		return "1.11", nil
	case standOrder63:
		return "1.12", nil
	default:
		return "", fmt.Errorf("stand %q: aiproj_version is required", standName)
	}
}

// isAIE5x reports whether the stand is Application Inspector Enterprise 5.x.
func isAIE5x(standName string) bool {
	v, err := StandVersion(standName)
	if err != nil {
		return false
	}

	return strings.HasPrefix(v, "5.")
}

func AiprojFixturePath(fixturesDir, aiprojVersion string) string {
	return filepath.Join(fixturesDir, fmt.Sprintf("aiproj-%s.json", aiprojVersion))
}

func OrderedStandNames(stands map[string]Stand) []string {
	names := make([]string, 0, len(standOrder))
	for _, name := range standOrder {
		if _, ok := stands[name]; ok {
			names = append(names, name)
		}
	}

	for name := range stands {
		found := false
		for _, ordered := range standOrder {
			if ordered == name {
				found = true
				break
			}
		}
		if !found {
			names = append(names, name)
		}
	}

	return names
}
