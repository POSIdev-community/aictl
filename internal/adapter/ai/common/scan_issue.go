package common

import "github.com/POSIdev-community/aictl/internal/core/domain/statistic"

// NewScanIssue maps API level/approval strings into the thin statistic domain issue.
func NewScanIssue(level, approvalState string) statistic.Issue {
	return statistic.Issue{
		Level:         statistic.Level(level),
		ApprovalState: statistic.ApprovalState(approvalState),
	}
}
