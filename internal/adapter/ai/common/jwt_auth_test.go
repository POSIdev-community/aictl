package common

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
)

func TestRefreshOrReauth_RefreshOK(t *testing.T) {
	t.Parallel()

	var reauthCalls atomic.Int32
	err := RefreshOrReauth(t.Context(), func(context.Context) error {
		return nil
	}, func(context.Context) error {
		reauthCalls.Add(1)

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, int32(0), reauthCalls.Load())
}

func TestRefreshOrReauth_RefreshNonAuthPropagates(t *testing.T) {
	t.Parallel()

	want := errors.New("network down")
	var reauthCalls atomic.Int32
	err := RefreshOrReauth(t.Context(), func(context.Context) error {
		return want
	}, func(context.Context) error {
		reauthCalls.Add(1)

		return nil
	})
	require.ErrorIs(t, err, want)
	require.Equal(t, int32(0), reauthCalls.Load())
}

func TestRefreshOrReauth_RefreshAuthThenReauthOK(t *testing.T) {
	t.Parallel()

	var reauthCalls atomic.Int32
	err := RefreshOrReauth(t.Context(), func(context.Context) error {
		return apperror.NewAuthenticationError()
	}, func(context.Context) error {
		reauthCalls.Add(1)

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, int32(1), reauthCalls.Load())
}

func TestRefreshOrReauth_Reauth401FailFast(t *testing.T) {
	t.Parallel()

	var reauthCalls atomic.Int32
	err := RefreshOrReauth(t.Context(), func(context.Context) error {
		return apperror.NewAuthenticationError()
	}, func(context.Context) error {
		reauthCalls.Add(1)

		return apperror.NewAuthenticationError()
	})
	require.Error(t, err)

	var authErr *apperror.AuthenticationError
	require.ErrorAs(t, err, &authErr)
	require.Equal(t, int32(1), reauthCalls.Load())
}

func TestRefreshOrReauth_Reauth5xxThenOK(t *testing.T) {
	t.Parallel()

	var reauthCalls atomic.Int32
	err := RefreshOrReauth(t.Context(), func(context.Context) error {
		return apperror.NewAuthenticationError()
	}, func(context.Context) error {
		if reauthCalls.Add(1) == 1 {
			return apperror.NewServerError(503, "unavailable")
		}

		return nil
	})
	require.NoError(t, err)
	require.Equal(t, int32(2), reauthCalls.Load())
}

func TestRefreshOrReauth_Reauth5xxRespectsCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	var reauthCalls atomic.Int32

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := RefreshOrReauth(ctx, func(context.Context) error {
		return apperror.NewAuthenticationError()
	}, func(context.Context) error {
		reauthCalls.Add(1)

		return apperror.NewServerError(500, "boom")
	})
	require.ErrorIs(t, err, context.Canceled)
	require.GreaterOrEqual(t, reauthCalls.Load(), int32(1))
}
