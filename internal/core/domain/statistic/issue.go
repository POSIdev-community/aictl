package statistic

// Level is the AI IssueLevel used for severity counters.
type Level string

const (
	LevelNone      Level = "None"
	LevelPotential Level = "Potential"
	LevelLow       Level = "Low"
	LevelMedium    Level = "Medium"
	LevelHigh      Level = "High"
)

// ApprovalState is the AI IssueApprovalState used for triage filtering.
type ApprovalState string

const (
	ApprovalNone         ApprovalState = "None"
	ApprovalApproval     ApprovalState = "Approval"
	ApprovalAutoApproval ApprovalState = "AutoApproval"
	ApprovalDiscard      ApprovalState = "Discard"
	ApprovalNotExist     ApprovalState = "NotExist"
)

// Issue is a thin scan finding used only to build severity statistics.
type Issue struct {
	Level         Level
	ApprovalState ApprovalState
}

// IsDiscarded reports whether triage discarded this finding.
func (i Issue) IsDiscarded() bool {
	return i.ApprovalState == ApprovalDiscard
}

// SeverityCounts is severity breakdown built from scan issues.
type SeverityCounts struct {
	High      int32
	Medium    int32
	Low       int32
	Potential int32
	// Other covers Level None / unknown — counted only in Total.
	Other int32
}

// Total returns the number of counted issues.
func (c SeverityCounts) Total() int32 {
	return c.High + c.Medium + c.Low + c.Potential + c.Other
}

func (c *SeverityCounts) addLevel(level Level) {
	switch level {
	case LevelHigh:
		c.High++
	case LevelMedium:
		c.Medium++
	case LevelLow:
		c.Low++
	case LevelPotential:
		c.Potential++
	default:
		c.Other++
	}
}

// CountSeverities builds severity counters from issues.
// When excludeDiscarded is true, Discard-triaged issues are skipped.
func CountSeverities(issues []Issue, excludeDiscarded bool) SeverityCounts {
	var counts SeverityCounts
	for _, issue := range issues {
		if excludeDiscarded && issue.IsDiscarded() {
			continue
		}
		counts.addLevel(issue.Level)
	}
	return counts
}

// ApplySeverityCounts overwrites vulnerability counters on s.
func (s *Statistic) ApplySeverityCounts(counts SeverityCounts) {
	if s == nil {
		return
	}
	s.High = counts.High
	s.Medium = counts.Medium
	s.Low = counts.Low
	s.Potential = counts.Potential
	s.Total = counts.Total()
}
