package v5_x

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/clientai/v5_x"
)

func TestMapScanIssues5x(t *testing.T) {
	t.Parallel()

	discard := v5_x.IssueApprovalStateDiscard
	none := v5_x.IssueApprovalStateNone
	high := v5_x.IssueLevelHigh
	medium := v5_x.IssueLevelMedium

	got := mapScanIssues5x([]v5_x.VulnerabilityModel{
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
