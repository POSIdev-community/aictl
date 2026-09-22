package v6_1

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_1"
)

func TestMapScanIssues61(t *testing.T) {
	t.Parallel()

	discard := v6_1.IssueApprovalStateDiscard
	none := v6_1.IssueApprovalStateNone
	high := v6_1.IssueLevelHigh
	medium := v6_1.IssueLevelMedium

	got := mapScanIssues61([]v6_1.VulnerabilityModel{
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
