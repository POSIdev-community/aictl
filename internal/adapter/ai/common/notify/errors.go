package notify

import (
	"errors"
	"fmt"
	"net/http"
)

// AuthError means the notification hub rejected the access token.
type AuthError struct {
	StatusCode int
	Body       string
}

func (e *AuthError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("notify auth failed: status %d", e.StatusCode)
	}

	return fmt.Sprintf("notify auth failed: status %d: %s", e.StatusCode, e.Body)
}

func IsAuthError(err error) bool {
	var authErr *AuthError

	return errors.As(err, &authErr)
}

func newAuthError(status int, body string) error {
	return &AuthError{StatusCode: status, Body: truncate(body, 256)}
}

func isUnauthorized(status int) bool {
	return status == http.StatusUnauthorized
}
