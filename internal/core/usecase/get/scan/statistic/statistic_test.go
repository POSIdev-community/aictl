package statistic

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
)

type fakeAI struct {
	stat   *statistic.Statistic
	issues []statistic.Issue
}

func (f *fakeAI) InitializeWithRetry(context.Context) error { return nil }

func (f *fakeAI) GetScanStatistic(context.Context, uuid.UUID, uuid.UUID) (*statistic.Statistic, error) {
	cp := *f.stat
	return &cp, nil
}

func (f *fakeAI) GetScanIssues(context.Context, uuid.UUID, uuid.UUID) ([]statistic.Issue, error) {
	return append([]statistic.Issue(nil), f.issues...), nil
}

type fakeCLI struct {
	logs   []string
	shown  *statistic.Statistic
	stdout bytes.Buffer
}

func (f *fakeCLI) ShowReader(r io.Reader) error {
	_, err := io.Copy(&f.stdout, r)
	return err
}

func (f *fakeCLI) ShowTextf(_ context.Context, format string, args ...any) {
	f.logs = append(f.logs, format)
	_ = args
}

func (f *fakeCLI) ShowScanStatistic(_ context.Context, s *statistic.Statistic) {
	f.shown = s
}

func TestUseCase_Execute_WithTriage(t *testing.T) {
	t.Parallel()

	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	scanID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	cfg := config.NewConfig(mustURI(t, "https://example.test"), "token", true, projectID, uuid.Nil)

	ai := &fakeAI{
		stat: &statistic.Statistic{
			FilesTotal:   10,
			FilesScanned: 8,
			PolicyState:  "None",
			Total:        99, // must be overwritten from issues
			High:         99,
		},
		issues: []statistic.Issue{
			{Level: statistic.LevelHigh, ApprovalState: statistic.ApprovalDiscard},
			{Level: statistic.LevelHigh, ApprovalState: statistic.ApprovalNone},
			{Level: statistic.LevelMedium, ApprovalState: statistic.ApprovalApproval},
			{Level: statistic.LevelLow, ApprovalState: statistic.ApprovalDiscard},
			{Level: statistic.LevelPotential, ApprovalState: statistic.ApprovalNone},
		},
	}

	t.Run("include_discard", func(t *testing.T) {
		cli := &fakeCLI{}
		uc, err := NewUseCase(ai, cli, cfg)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(t.Context(), scanID, "", false, false))
		require.NotNil(t, cli.shown)
		require.Equal(t, int32(2), cli.shown.High)
		require.Equal(t, int32(1), cli.shown.Medium)
		require.Equal(t, int32(1), cli.shown.Low)
		require.Equal(t, int32(1), cli.shown.Potential)
		require.Equal(t, int32(5), cli.shown.Total)
		require.Equal(t, int32(10), cli.shown.FilesTotal)
		require.Contains(t, cli.logs, "getting scan statistic, scan-id '%v'")
		require.Contains(t, cli.logs, "getting scan issues, scan-id '%v'")
	})

	t.Run("exclude_discard", func(t *testing.T) {
		cli := &fakeCLI{}
		uc, err := NewUseCase(ai, cli, cfg)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(t.Context(), scanID, "", false, true))
		require.NotNil(t, cli.shown)
		require.Equal(t, int32(1), cli.shown.High)
		require.Equal(t, int32(1), cli.shown.Medium)
		require.Equal(t, int32(0), cli.shown.Low)
		require.Equal(t, int32(1), cli.shown.Potential)
		require.Equal(t, int32(3), cli.shown.Total)
	})
}

func mustURI(t *testing.T, raw string) config.Uri {
	t.Helper()
	uri, err := config.NewUri(raw)
	require.NoError(t, err)
	return uri
}
