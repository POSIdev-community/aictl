package notify

import (
	"context"
	"math/rand/v2"
	"time"
)

const (
	ReconnectMinDelay               = time.Second
	ReconnectMaxDelay               = 30 * time.Second
	StableConnectionForBackoffReset = 30 * time.Second
)

// NextBackoff doubles delay up to ReconnectMaxDelay.
func NextBackoff(current time.Duration) time.Duration {
	if current < ReconnectMinDelay {
		return ReconnectMinDelay
	}

	next := current * 2
	if next > ReconnectMaxDelay {
		return ReconnectMaxDelay
	}

	return next
}

// Sleep waits for d with ±20% jitter, or until ctx is cancelled.
// Returns false if ctx was cancelled.
func Sleep(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		d = ReconnectMinDelay
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
