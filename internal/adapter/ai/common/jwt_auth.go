package common

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
)

const (
	authRetryMinDelay = time.Second
	authRetryMaxDelay = 30 * time.Second
)

// RefreshOrReauth tries refresh first. On AuthenticationError it falls back to reauth
// (sign-in with API token). Reauth AuthenticationError fails immediately; ServerResponseError
// (5xx) is retried with exponential backoff until ctx is cancelled.
func RefreshOrReauth(ctx context.Context, refresh, reauth func(context.Context) error) error {
	if err := refresh(ctx); err == nil {
		return nil
	} else if !isAuthenticationError(err) {
		return err
	}

	return reauthWithServerRetry(ctx, reauth)
}

func reauthWithServerRetry(ctx context.Context, reauth func(context.Context) error) error {
	backoff := authRetryMinDelay

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := reauth(ctx)
		if err == nil {
			return nil
		}
		if isAuthenticationError(err) {
			return err
		}
		if !isServerResponseError(err) {
			return err
		}

		if !sleepAuthRetry(ctx, backoff) {
			return ctx.Err()
		}
		backoff = nextAuthBackoff(backoff)
	}
}

func isAuthenticationError(err error) bool {
	var authErr *apperror.AuthenticationError

	return errors.As(err, &authErr)
}

func isServerResponseError(err error) bool {
	var serverErr *apperror.ServerResponseError

	return errors.As(err, &serverErr)
}

func nextAuthBackoff(current time.Duration) time.Duration {
	if current < authRetryMinDelay {
		return authRetryMinDelay
	}

	next := current * 2
	if next > authRetryMaxDelay {
		return authRetryMaxDelay
	}

	return next
}

func sleepAuthRetry(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		d = authRetryMinDelay
	}

	jitter := time.Duration(float64(d) * (0.8 + 0.4*rand.Float64()))
	timer := time.NewTimer(jitter)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
