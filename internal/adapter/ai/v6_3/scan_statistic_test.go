package v6_3

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_3"
)

func TestMapScanIssues6x(t *testing.T) {
	t.Parallel()

	discard := v6_3.IssueApprovalStateDiscard
	none := v6_3.IssueApprovalStateNone
	high := v6_3.IssueLevelHigh
	medium := v6_3.IssueLevelMedium

	got := mapScanIssues6x([]v6_3.VulnerabilityModel{
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
