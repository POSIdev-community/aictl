package v6_4

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_4"
)

func TestMapScanIssues6x(t *testing.T) {
	t.Parallel()

	discard := v6_4.IssueApprovalStateDiscard
	none := v6_4.IssueApprovalStateNone
	high := v6_4.IssueLevelHigh
	medium := v6_4.IssueLevelMedium

	got := mapScanIssues6x([]v6_4.VulnerabilityModel{
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
