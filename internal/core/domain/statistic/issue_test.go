package statistic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountSeverities(t *testing.T) {
	t.Parallel()

	issues := []Issue{
		{Level: LevelHigh, ApprovalState: ApprovalDiscard},
		{Level: LevelHigh, ApprovalState: ApprovalNone},
		{Level: LevelMedium, ApprovalState: ApprovalApproval},
		{Level: LevelLow, ApprovalState: ApprovalDiscard},
		{Level: LevelPotential, ApprovalState: ApprovalNone},
		{Level: LevelNone, ApprovalState: ApprovalNone},
	}

	all := CountSeverities(issues, false)
	require.Equal(t, int32(2), all.High)
	require.Equal(t, int32(1), all.Medium)
	require.Equal(t, int32(1), all.Low)
	require.Equal(t, int32(1), all.Potential)
	require.Equal(t, int32(1), all.Other)
	require.Equal(t, int32(6), all.Total())

	triaged := CountSeverities(issues, true)
	require.Equal(t, int32(1), triaged.High)
	require.Equal(t, int32(1), triaged.Medium)
	require.Equal(t, int32(0), triaged.Low)
	require.Equal(t, int32(1), triaged.Potential)
	require.Equal(t, int32(1), triaged.Other)
	require.Equal(t, int32(4), triaged.Total())
}

func TestStatistic_ApplySeverityCounts(t *testing.T) {
	t.Parallel()

	stat := &Statistic{
		Total:        99,
		High:         99,
		FilesTotal:   10,
		FilesScanned: 5,
		PolicyState:  "None",
	}
	stat.ApplySeverityCounts(SeverityCounts{High: 1, Medium: 2, Low: 3, Potential: 4})

	require.Equal(t, int32(1), stat.High)
	require.Equal(t, int32(2), stat.Medium)
	require.Equal(t, int32(3), stat.Low)
	require.Equal(t, int32(4), stat.Potential)
	require.Equal(t, int32(10), stat.Total)
	require.Equal(t, int32(10), stat.FilesTotal)
	require.Equal(t, "None", stat.PolicyState)

	(*Statistic)(nil).ApplySeverityCounts(SeverityCounts{High: 1})
}
