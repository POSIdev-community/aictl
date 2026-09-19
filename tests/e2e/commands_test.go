//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestLeafCommandsOutsideSmoke exercises implemented leaf commands that are not
// covered by run-pipeline.sh. Asserts exit 0 across configured AIE stands.
//
// Skipped (no safe happy-path without extra fixtures / version support):
//   - get scan report <custom-template> — needs a real template name on the server
//   - get scan sbom — e2e aiproj has only StaticCodeAnalysis (API returns sbom not found)
//   - get/update project settings — not on AIE 5.x
//   - get scan report json-v2 — AIE 6.1+ only
//   - get scan report xml — on AIE 6.1+ expect error (template removed)
func TestLeafCommandsOutsideSmoke(t *testing.T) {
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

	root, err := RepoRoot()
	require.NoError(t, err)

	aictlBin, err := ResolveAictlBin(root)
	require.NoError(t, err, "build aictl for e2e")
	t.Logf("using aictl binary: %s", aictlBin)

	fixturesDir := filepath.Join(root, "tests", "e2e", "fixtures")

	for _, standName := range OrderedStandNames(stands) {
		stand := stands[standName]

		t.Run("AIE_"+standName, func(t *testing.T) {
			t.Parallel()

			aiprojVersion, err := stand.ResolveAiprojVersion(standName)
			require.NoError(t, err)

			workDir := t.TempDir()
			xdgConfig := filepath.Join(workDir, "xdg-config")
			require.NoError(t, os.MkdirAll(xdgConfig, 0o755))
			env := []string{"XDG_CONFIG_HOME=" + xdgConfig}

			projectName := fmt.Sprintf("aictl-e2e-leaf-%s-%s", standName, uuid.NewString())
			aiprojPath := filepath.Join(workDir, "aiproj.json")
			writePatchedAiproj(t, AiprojFixturePath(fixturesDir, aiprojVersion), aiprojPath, projectName)

			run := func(args ...string) string {
				return RunAictl(t, aictlBin, stand, env, args...)
			}

			projectID := run("create", "project", projectName, "--safe")
			assertUUID(t, projectID)
			t.Cleanup(func() {
				_, _ = runAictl(t, aictlBin, stand, env, "delete", "projects", projectID)
			})

			run("set", "project", "settings", "-p", projectID, "-f", aiprojPath)
			branchID := run("create", "branch", "default", "-p", projectID, "--safe")
			assertUUID(t, branchID)
			run("update", "sources", filepath.Join(fixturesDir, "project"), "-p", projectID, "-b", branchID, "--update-languages")

			// stop: start a branch scan and stop it while still running.
			// Await the stopped scan before starting another — otherwise AIE may
			// still treat the branch as scheduled (SCAN_ALREADY_SCHEDULED flake).
			stopScanID := run("scan", "branch", branchID, "-p", projectID)
			assertUUID(t, stopScanID)
			run("scan", "stop", stopScanID)
			run("scan", "await", stopScanID, "-p", projectID)

			scanID := run("scan", "branch", branchID, "-p", projectID)
			assertUUID(t, scanID)
			run("scan", "await", scanID, "-p", projectID)

			reportsDir := filepath.Join(workDir, "reports")
			require.NoError(t, os.MkdirAll(reportsDir, 0o755))

			t.Run("ctx", func(t *testing.T) {
				RunAictlNoConn(t, aictlBin, env,
					"ctx", "set",
					"-u", stand.URL,
					"-t", stand.Token,
					"--tls-skip",
					"-p", projectID,
					"-b", branchID,
				)
				RunAictlNoConn(t, aictlBin, env, "ctx", "show")
				RunAictlNoConn(t, aictlBin, env, "ctx", "unset", "-p", "-b")
				RunAictlNoConn(t, aictlBin, env, "ctx", "clear", "-y")
			})

			commands := []struct {
				name string
				args []string
			}{
				{"get healthcheck", []string{"get", "healthcheck"}},
				{"get agents", []string{"get", "agents"}},
				{"get projects", []string{"get", "projects", regexpEscape(projectName)}},
				{"get project aiproj", []string{"get", "project", "aiproj", "-p", projectID}},
				{"get project policies", []string{"get", "project", "policies", "-p", projectID}},
				{"get project exclusions", []string{"get", "project", "exclusions", "-p", projectID}},
				{"get branches", []string{"get", "branches", "-p", projectID}},
				{"get branch", []string{"get", "branch", branchID}},
				{"get scans", []string{"get", "scans", "-b", branchID}},
				{"get scan", []string{"get", "scan", scanID, "-p", projectID}},
				{"get scan aiproj", []string{"get", "scan", "aiproj", scanID, "-p", projectID}},
				{"get scan logs", []string{"get", "scan", "logs", scanID, "-p", projectID, "-o", filepath.Join(reportsDir, "logs.txt")}},
				{"get scan errors", []string{"get", "scan", "errors", scanID, "-p", projectID}},
				// get scan sbom omitted: fixture ScanModules is StaticCodeAnalysis only → "sbom not found"
				{"get scan stage", []string{"get", "scan", "stage", scanID, "-p", projectID}},
				{"get scan stage --fail-on-scan-failed", []string{"get", "scan", "stage", scanID, "-p", projectID, "--fail-on-scan-failed"}},
				{"get scan statistic", []string{"get", "scan", "statistic", scanID, "-p", projectID}},
				{"get queue", []string{"get", "queue"}},
				{"get queue -p", []string{"get", "queue", "-p", projectID}},
				{"get scanning", []string{"get", "scanning"}},
				{"get scanning -p", []string{"get", "scanning", "-p", projectID}},
				{"get report-templates", []string{"get", "report-templates", "--localization", "en"}},

				{"scan check-policies", []string{"scan", "check-policies", scanID, "-p", projectID}},
				{"scan await --fail-on-scan-failed", []string{"scan", "await", scanID, "-p", projectID, "--fail-on-scan-failed"}},
			}

			// priority / preferred-agents / languages settings exist only on AIE 6.0+
			if !isAIE5x(standName) {
				commands = append(commands,
					struct {
						name string
						args []string
					}{"get project settings", []string{"get", "project", "settings", "-p", projectID}},
					struct {
						name string
						args []string
					}{"update project settings", []string{"update", "project", "settings", "-p", projectID, "--priority", "Medium"}},
					struct {
						name string
						args []string
					}{"update project languages", []string{"update", "project", "languages", "-p", projectID}},
				)
			}

			// SCA feeds packages API exists only on AIE 6.3+
			if standName == standOrder63 {
				commands = append(commands,
					struct {
						name string
						args []string
					}{"get sca-feeds", []string{"get", "sca-feeds"}},
					struct {
						name string
						args []string
					}{"get sca-feeds --status current", []string{"get", "sca-feeds", "--status", "current"}},
				)
			}

			reportTypes := []string{
				"autocheck", "gitlab", "json", "markdown",
				"nist", "oud4", "owasp", "owaspm", "pcidss", "plain", "sans", "sarif", "xml",
			}
			// json-v2 is supported on AIE 6.1+ only
			if standName == standOrder61 || standName == standOrder62 || standName == standOrder63 {
				reportTypes = append(reportTypes, "json-v2")
			}
			xmlExpectError := standName == standOrder61 || standName == standOrder62 || standName == standOrder63
			for _, rt := range reportTypes {
				out := filepath.Join(reportsDir, rt+".out")
				commands = append(commands, struct {
					name string
					args []string
				}{
					name: "get scan report " + rt,
					args: []string{"get", "scan", "report", rt, scanID, "-p", projectID, "-o", out, "--localization", "en"},
				})
			}

			uiFilterFlags := uiLikeReportFilterFlags(standName)
			for _, rt := range reportTypes {
				out := filepath.Join(reportsDir, "filtered-"+rt+".out")
				args := []string{"get", "scan", "report", "with-filters", rt, scanID, "-p", projectID, "-o", out, "--localization", "en"}
				args = append(args, uiFilterFlags...)
				commands = append(commands, struct {
					name string
					args []string
				}{
					name: "get scan report with-filters " + rt,
					args: args,
				})
			}

			for _, c := range commands {
				c := c
				t.Run(c.name, func(t *testing.T) {
					if xmlExpectError && isXmlReportCommand(c.args) {
						_, err := runAictl(t, aictlBin, stand, env, c.args...)
						require.Error(t, err)
						return
					}
					RunAictl(t, aictlBin, stand, env, c.args...)
				})
			}

			t.Run("get scan report with-filters without filters", func(t *testing.T) {
				_, err := runAictl(t, aictlBin, stand, env,
					"get", "scan", "report", "with-filters", "sarif", scanID, "-p", projectID,
					"-o", filepath.Join(reportsDir, "no-filters.sarif"), "--localization", "en")
				require.Error(t, err)
			})

			if isAIE5x(standName) {
				t.Run("get scan report with-filters unsupported SecretDetection on 5.x", func(t *testing.T) {
					args := []string{
						"get", "scan", "report", "with-filters", "sarif", scanID, "-p", projectID,
						"-o", filepath.Join(reportsDir, "secret-detection-5x.sarif"),
						"--localization", "en",
						"--level-high",
						"--scan-module", "SecretDetection",
					}
					_, err := runAictl(t, aictlBin, stand, env, args...)
					require.Error(t, err)
				})
			}

			t.Run("set project policies", func(t *testing.T) {
				// API: securityPolicies is a JSON *string* (not an array).
				policiesPath := filepath.Join(workDir, "policies.json")
				require.NoError(t, os.WriteFile(policiesPath, []byte(
					`{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[]"}`,
				), 0o644))
				RunAictl(t, aictlBin, stand, env, "set", "project", "policies", "-p", projectID, "-f", policiesPath)
			})

			t.Run("set project exclusions", func(t *testing.T) {
				exclusionsPath := filepath.Join(workDir, "exclusions.gitignore")
				require.NoError(t, os.WriteFile(exclusionsPath, []byte("*.log\nnode_modules/\n"), 0o644))
				RunAictl(t, aictlBin, stand, env, "set", "project", "exclusions", "-p", projectID, "-f", exclusionsPath)
				out := RunAictl(t, aictlBin, stand, env, "get", "project", "exclusions", "-p", projectID)
				require.Contains(t, out, "*.log")
			})

			t.Run("update sources --temp-dir", func(t *testing.T) {
				tempDir := filepath.Join(workDir, "zip-temp")
				require.NoError(t, os.MkdirAll(tempDir, 0o755))
				RunAictl(t, aictlBin, stand, env,
					"update", "sources", filepath.Join(fixturesDir, "project"),
					"-p", projectID, "-b", branchID, "--temp-dir", tempDir, "--update-languages")
			})

			t.Run("scan check-policies --fail-on-policies-rejected", func(t *testing.T) {
				state := RunAictl(t, aictlBin, stand, env, "scan", "check-policies", scanID, "-p", projectID)
				if strings.TrimSpace(state) == "Rejected" {
					t.Skip("PolicyState is Rejected; --fail-on-policies-rejected would exit 1")
				}
				RunAictl(t, aictlBin, stand, env,
					"scan", "check-policies", scanID, "-p", projectID, "--fail-on-policies-rejected")
			})

			// start project last: exit 0 only (stop/await may hit SCAN*_NOT_FOUND on some AIE).
			t.Run("scan project", func(t *testing.T) {
				projectScanID := RunAictl(t, aictlBin, stand, env, "scan", "project", projectID)
				assertUUID(t, projectScanID)
			})
		})
	}
}

func writePatchedAiproj(t *testing.T, src, dst, projectName string) {
	t.Helper()

	data, err := os.ReadFile(src)
	require.NoError(t, err)

	var doc map[string]any
	require.NoError(t, json.Unmarshal(data, &doc))
	doc["ProjectName"] = projectName

	out, err := json.Marshal(doc)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(dst, out, 0o644))
}

// uiLikeReportFilterFlags returns UI-default-like filter flags for report with-filters e2e.
// SecretDetection / MaliciousCodeDetection are included only on AIE ≥ 6.0.
func uiLikeReportFilterFlags(standName string) []string {
	flags := []string{
		"--level-high", "--level-medium",
		"--status-undefined", "--status-confirmed", "--status-confirmed-auto",
		"--mode-entry-point", "--mode-root-function", "--mode-public-methods", "--mode-others",
		"--found-this-scan", "--found-prev-scan",
		"--conditional", "--non-conditional",
		"--non-suppressed",
		"--suspected", "--second-level",
		"--scan-module", "StaticCodeAnalysis",
		"--scan-module", "PatternMatching",
		"--scan-module", "Components",
		"--scan-module", "SoftwareCompositionAnalysis",
		"--scan-module", "Configuration",
		"--scan-module", "BlackBox",
	}
	if !isAIE5x(standName) {
		flags = append(flags,
			"--scan-module", "MaliciousCodeDetection",
			"--scan-module", "SecretDetection",
		)
	}
	return flags
}

// isXmlReportCommand reports whether args are get scan report [with-filters] xml …
func isXmlReportCommand(args []string) bool {
	if len(args) < 4 {
		return false
	}
	if args[0] != "get" || args[1] != "scan" || args[2] != "report" {
		return false
	}
	if args[3] == "xml" {
		return true
	}
	return len(args) >= 5 && args[3] == "with-filters" && args[4] == "xml"
}

func regexpEscape(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`.`, `\.`,
		`+`, `\+`,
		`*`, `\*`,
		`?`, `\?`,
		`(`, `\(`,
		`)`, `\)`,
		`[`, `\[`,
		`]`, `\]`,
		`{`, `\{`,
		`}`, `\}`,
		`|`, `\|`,
		`^`, `\^`,
		`$`, `\$`,
	)

	return replacer.Replace(s)
}
