package apperror

import "fmt"

// FailError is returned when an opt-in fail flag triggers (exit code 1).
// It is distinct from validation.Error.
type FailError struct {
	message string
}

func (e *FailError) Error() string {
	return e.message
}

func NewFailError(message string) *FailError {
	return &FailError{message: message}
}

func NewScanFailedError(stage string) *FailError {
	return &FailError{message: fmt.Sprintf("scan failed with stage %s", stage)}
}

func NewPolicyFailError(state string) *FailError {
	return &FailError{message: fmt.Sprintf("policy check failed: %s", state)}
}
