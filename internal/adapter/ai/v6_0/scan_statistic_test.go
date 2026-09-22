package v6_0

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_0"
)

func TestMapScanIssues60(t *testing.T) {
	t.Parallel()

	discard := v6_0.IssueApprovalStateDiscard
	none := v6_0.IssueApprovalStateNone
	high := v6_0.IssueLevelHigh
	medium := v6_0.IssueLevelMedium

	got := mapScanIssues60([]v6_0.VulnerabilityModel{
		{ApprovalState: &discard, Level: &high},
		{ApprovalState: &none, Level: &medium},
		{ApprovalState: nil, Level: &high},
	})

	require.Equal(t, []statistic.Issue{
		{Level: statistic.LevelHigh, ApprovalState: statistic.ApprovalDiscard},
		{Level: statistic.LevelMedium, ApprovalState: statistic.ApprovalNone},
		{Level: statistic.LevelHigh, ApprovalState: ""},
	}, got)
}
