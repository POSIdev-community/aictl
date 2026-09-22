package v6_x

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
)

func TestMapScanIssues6x(t *testing.T) {
	t.Parallel()

	discard := v6_x.IssueApprovalStateDiscard
	none := v6_x.IssueApprovalStateNone
	high := v6_x.IssueLevelHigh
	medium := v6_x.IssueLevelMedium

	got := mapScanIssues6x([]v6_x.VulnerabilityModel{
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
