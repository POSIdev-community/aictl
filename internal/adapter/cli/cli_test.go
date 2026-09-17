package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

func newTestLogger(out, err *bytes.Buffer) *zap.Logger {
	encCfg := zapcore.EncoderConfig{
		MessageKey: "msg",
		LineEnding: zapcore.DefaultLineEnding,
	}
	infoCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(out),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		}),
	)
	errCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(err),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		}),
	)
	return zap.New(zapcore.NewTee(infoCore, errCore))
}

func testCtx(out, errBuf *bytes.Buffer) context.Context {
	return logger.ContextWithLogger(context.Background(), logger.Wrap(newTestLogger(out, errBuf)))
}

func newTestAdapter(in io.Reader, out io.Writer) *Adapter {
	if in == nil {
		in = strings.NewReader("")
	}
	if out == nil {
		out = io.Discard
	}
	return &Adapter{stdin: in, stdout: out}
}

func TestReturnAndShowText(t *testing.T) {
	var out, errBuf bytes.Buffer
	ctx := testCtx(&out, &errBuf)
	a := newTestAdapter(nil, nil)

	a.ReturnText(ctx, "result")
	a.ReturnTextf(ctx, "fmt-%s", "x")
	a.ShowText(ctx, "msg")
	a.ShowTextf(ctx, "err-%s", "y")

	require.Contains(t, out.String(), "result")
	require.Contains(t, out.String(), "fmt-x")
	require.Contains(t, errBuf.String(), "msg")
	require.Contains(t, errBuf.String(), "err-y")
}

func TestShowProjects(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	projects := []project.Project{project.NewProject(id, "demo")}

	t.Run("table", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowProjects(testCtx(&out, &errBuf), projects)
		require.Contains(t, out.String(), "ID")
		require.Contains(t, out.String(), "NAME")
		require.Contains(t, out.String(), "TYPE")
		require.Contains(t, out.String(), id.String())
		require.Contains(t, out.String(), "demo")
		require.Contains(t, out.String(), "source")
	})

	t.Run("quiet", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowProjectsQuite(testCtx(&out, &errBuf), projects)
		require.Contains(t, out.String(), id.String())
		require.NotContains(t, out.String(), "NAME")
	})
}

func TestShowBranches(t *testing.T) {
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	branches := []branch.Branch{branch.NewBranch(id, "main", "", true)}

	t.Run("table", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowBranches(testCtx(&out, &errBuf), branches)
		require.Contains(t, out.String(), id.String())
		require.Contains(t, out.String(), "main")
	})

	t.Run("quiet", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowBranchesQuite(testCtx(&out, &errBuf), branches)
		require.Contains(t, out.String(), id.String())
	})
}

func TestShowScans(t *testing.T) {
	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	date := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	label := "nightly"
	scans := []scan.Scan{scan.NewScan(id, uuid.Nil, &date, &label)}

	t.Run("table", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowScans(testCtx(&out, &errBuf), scans)
		require.Contains(t, out.String(), id.String())
		require.Contains(t, out.String(), "2024-01-02 03:04:05")
		require.Contains(t, out.String(), "nightly")
	})

	t.Run("quiet", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowScansQuite(testCtx(&out, &errBuf), scans)
		require.Contains(t, out.String(), id.String())
	})
}

func TestShowScanAgents(t *testing.T) {
	id := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	agents := []scanagent.ScanAgent{{
		Id: id, Name: "agent-1", Status: "Online", Version: "1.0", OperatingSystem: "linux",
	}}

	t.Run("table", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowScanAgents(testCtx(&out, &errBuf), agents)
		require.Contains(t, out.String(), id.String())
		require.Contains(t, out.String(), "agent-1")
		require.Contains(t, out.String(), "Online")
	})

	t.Run("quiet", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowScanAgentsQuite(testCtx(&out, &errBuf), agents)
		require.Contains(t, out.String(), id.String())
	})
}

func TestShowQueue(t *testing.T) {
	var out, errBuf bytes.Buffer
	a := newTestAdapter(nil, nil)
	pid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	bid := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	sid := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	a.ShowQueue(testCtx(&out, &errBuf), []queue.Entry{{
		ProjectId: pid, BranchId: bid, ScanId: sid, Stage: "Scanning",
	}})

	s := out.String()
	require.Contains(t, s, pid.String())
	require.Contains(t, s, bid.String())
	require.Contains(t, s, sid.String())
	require.Contains(t, s, "Scanning")
}

func TestShowReportTemplates(t *testing.T) {
	id := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	templates := []report.Template{{Id: id, Name: "SARIF"}}

	t.Run("table", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowReportTemplates(testCtx(&out, &errBuf), templates)
		require.Contains(t, out.String(), id.String())
		require.Contains(t, out.String(), "SARIF")
	})

	t.Run("quiet", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowReportTemplatesQuite(testCtx(&out, &errBuf), templates)
		require.Contains(t, out.String(), id.String())
	})
}

func TestShowScanStatistic(t *testing.T) {
	var out, errBuf bytes.Buffer
	a := newTestAdapter(nil, nil)
	a.ShowScanStatistic(testCtx(&out, &errBuf), &statistic.Statistic{
		Total: 10, High: 4, Medium: 3, Low: 2, Potential: 1,
	})
	s := out.String()
	require.Contains(t, s, "Total: 10")
	require.Contains(t, s, "High: 4")
	require.Contains(t, s, "Medium: 3")
	require.Contains(t, s, "Low: 2")
	require.Contains(t, s, "Potential: 1")
}

func TestShowProjectSettings(t *testing.T) {
	agentID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	t.Run("with_agents", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowProjectSettings(testCtx(&out, &errBuf), settings.ProjectSettingsView{
			Priority:            settings.PriorityHigh,
			PreferredAgentsOnly: true,
			PreferredAgents:     []uuid.UUID{agentID},
		})
		s := out.String()
		require.Contains(t, s, "Priority: High")
		require.Contains(t, s, "Preferred agents only: true")
		require.Contains(t, s, agentID.String())
	})

	t.Run("none", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(nil, nil)
		a.ShowProjectSettings(testCtx(&out, &errBuf), settings.ProjectSettingsView{
			Priority: settings.PriorityLow,
		})
		require.Contains(t, out.String(), "Preferred agents: none")
	})
}

func TestAskConfirmation(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(strings.NewReader("yes\n"), nil)
		ok, err := a.AskConfirmation(testCtx(&out, &errBuf), "continue?")
		require.NoError(t, err)
		require.True(t, ok)
		require.Contains(t, out.String(), "continue?")
	})

	t.Run("no", func(t *testing.T) {
		var out, errBuf bytes.Buffer
		a := newTestAdapter(strings.NewReader("n\n"), nil)
		ok, err := a.AskConfirmation(testCtx(&out, &errBuf), "continue?")
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestShowReader(t *testing.T) {
	var stdout bytes.Buffer
	a := newTestAdapter(nil, &stdout)
	require.NoError(t, a.ShowReader(strings.NewReader("report-body")))
	require.Contains(t, stdout.String(), "report-body")
}
