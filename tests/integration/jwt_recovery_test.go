package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_x"
	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type tokenClient interface {
	Initialize(ctx context.Context, cfg *config.Config) error
	RefreshAccessToken(ctx context.Context) error
	GetAccessToken() string
	GetRefreshToken() string
}

func TestTokenRecovery_V5x(t *testing.T) {
	t.Parallel()
	runTokenRecoveryTests(t, func(base *common.BaseClient) tokenClient {
		return v5_x.NewAiClient(base)
	})
}

func TestTokenRecovery_V60(t *testing.T) {
	t.Parallel()
	runTokenRecoveryTests(t, func(base *common.BaseClient) tokenClient {
		return v6_0.NewAiClient(base)
	})
}

func TestTokenRecovery_V6x(t *testing.T) {
	t.Parallel()
	runTokenRecoveryTests(t, func(base *common.BaseClient) tokenClient {
		return v6_x.NewAiClient(base)
	})
}

func runTokenRecoveryTests(t *testing.T, newClient func(base *common.BaseClient) tokenClient) {
	t.Helper()

	const apiToken = "api-token-for-tests"

	t.Run("signin_401_at_initialize", func(t *testing.T) {
		t.Parallel()

		srv := NewAuthServer(t, apiToken, []Step{Unauthorized()}, nil)
		client := newClient(common.NewBaseClient())
		cfg := mustConfig(t, srv.URL, apiToken)

		err := client.Initialize(t.Context(), cfg)
		require.Error(t, err)

		var authErr *apperror.AuthenticationError
		require.ErrorAs(t, err, &authErr)
		require.Equal(t, 1, srv.SigninCalls)
		require.Equal(t, apiToken, srv.LastAPIToken)
		require.Equal(t, 0, srv.RefreshCalls)
	})

	t.Run("refresh_ok_continues", func(t *testing.T) {
		t.Parallel()

		srv := NewAuthServer(t, apiToken, []Step{
			OKSignin("access-1", "refresh-1"),
		}, []Step{
			OKRefresh("access-2"),
		})
		client := newClient(common.NewBaseClient())
		cfg := mustConfig(t, srv.URL, apiToken)

		require.NoError(t, client.Initialize(t.Context(), cfg))
		require.Equal(t, "access-1", client.GetAccessToken())
		require.Equal(t, "refresh-1", client.GetRefreshToken())

		require.NoError(t, client.RefreshAccessToken(t.Context()))
		require.Equal(t, "access-2", client.GetAccessToken())
		require.Equal(t, "refresh-1", client.GetRefreshToken())
		require.Equal(t, 1, srv.SigninCalls)
		require.Equal(t, 1, srv.RefreshCalls)
	})

	t.Run("refresh_401_then_reauth_ok", func(t *testing.T) {
		t.Parallel()

		srv := NewAuthServer(t, apiToken, []Step{
			OKSignin("access-1", "refresh-1"),
			OKSignin("access-3", "refresh-3"),
		}, []Step{
			Unauthorized(),
		})
		client := newClient(common.NewBaseClient())
		cfg := mustConfig(t, srv.URL, apiToken)

		require.NoError(t, client.Initialize(t.Context(), cfg))
		require.NoError(t, client.RefreshAccessToken(t.Context()))
		require.Equal(t, "access-3", client.GetAccessToken())
		require.Equal(t, "refresh-3", client.GetRefreshToken())
		require.Equal(t, 2, srv.SigninCalls)
		require.Equal(t, 1, srv.RefreshCalls)
		require.Equal(t, apiToken, srv.LastAPIToken)
	})

	t.Run("refresh_401_then_reauth_401_fail_fast", func(t *testing.T) {
		t.Parallel()

		srv := NewAuthServer(t, apiToken, []Step{
			OKSignin("access-1", "refresh-1"),
			Unauthorized(),
		}, []Step{
			Unauthorized(),
		})
		client := newClient(common.NewBaseClient())
		cfg := mustConfig(t, srv.URL, apiToken)

		require.NoError(t, client.Initialize(t.Context(), cfg))
		err := client.RefreshAccessToken(t.Context())
		require.Error(t, err)

		var authErr *apperror.AuthenticationError
		require.ErrorAs(t, err, &authErr)
		require.Equal(t, 2, srv.SigninCalls)
		require.Equal(t, 1, srv.RefreshCalls)
	})

	t.Run("refresh_401_then_reauth_5xx_then_ok", func(t *testing.T) {
		t.Parallel()

		srv := NewAuthServer(t, apiToken, []Step{
			OKSignin("access-1", "refresh-1"),
			ServerError(http.StatusServiceUnavailable),
			OKSignin("access-4", "refresh-4"),
		}, []Step{
			Unauthorized(),
		})
		client := newClient(common.NewBaseClient())
		cfg := mustConfig(t, srv.URL, apiToken)

		require.NoError(t, client.Initialize(t.Context(), cfg))
		require.NoError(t, client.RefreshAccessToken(t.Context()))
		require.Equal(t, "access-4", client.GetAccessToken())
		require.Equal(t, "refresh-4", client.GetRefreshToken())
		require.Equal(t, 3, srv.SigninCalls)
		require.Equal(t, 1, srv.RefreshCalls)
	})
}

func mustConfig(t *testing.T, rawURL, token string) *config.Config {
	t.Helper()

	uri, err := config.NewUri(rawURL)
	require.NoError(t, err)

	return config.NewConfig(uri, token, true, uuid.Nil, uuid.Nil)
}
