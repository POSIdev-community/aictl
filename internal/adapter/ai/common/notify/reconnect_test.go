package notify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNextBackoff(t *testing.T) {
	t.Parallel()

	require.Equal(t, time.Second, NextBackoff(0))
	require.Equal(t, 2*time.Second, NextBackoff(time.Second))
	require.Equal(t, 4*time.Second, NextBackoff(2*time.Second))
	require.Equal(t, ReconnectMaxDelay, NextBackoff(ReconnectMaxDelay))
	require.Equal(t, ReconnectMaxDelay, NextBackoff(20*time.Second))
}

func TestIsAuthError(t *testing.T) {
	t.Parallel()

	require.True(t, IsAuthError(newAuthError(401, "unauthorized")))
	require.False(t, IsAuthError(contextCanceledLike{}))
}

type contextCanceledLike struct{}

func (contextCanceledLike) Error() string { return "canceled" }
