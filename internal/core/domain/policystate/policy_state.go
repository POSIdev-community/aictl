package policystate

// State is the scan PolicyState enum from the AI API.
type State string

const (
	None      State = "None"
	Rejected  State = "Rejected"
	Confirmed State = "Confirmed"
)

func (s State) String() string {
	return string(s)
}

// IsRejected reports whether --fail-on-policies-rejected should exit non-zero.
func (s State) IsRejected() bool {
	return s == Rejected
}
