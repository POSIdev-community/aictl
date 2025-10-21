package clientai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryRoundTripper_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("should retry on error", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(
			func() http.Handler {
				var queryCount uint8

				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if queryCount == 0 {
						queryCount++
						w.WriteHeader(http.StatusUnauthorized)

						return
					}

					w.WriteHeader(http.StatusOK)
				})
			}(),
		)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error { close(retryCh); return nil })

		res, err := rtt.RoundTrip(httptest.NewRequest(http.MethodGet, srv.URL, nil))
		require.NoError(t, err)

		select {
		case <-retryCh:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for retry")
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)

	})

	t.Run("should retry with error", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(
			func() http.Handler {
				var queryCount uint8

				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if queryCount == 0 {
						queryCount++
						w.WriteHeader(http.StatusUnauthorized)

						return
					}

					w.WriteHeader(http.StatusOK)
				})
			}(),
		)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error { close(retryCh); return http.ErrHandlerTimeout })

		res, err := rtt.RoundTrip(httptest.NewRequest(http.MethodGet, srv.URL, nil))
		require.Error(t, err)
		assert.ErrorIs(t, err, http.ErrHandlerTimeout)

		select {
		case <-retryCh:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for retry")
		}

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	})

	t.Run("should not retry on success", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(
			func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
			}(),
		)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error { close(retryCh); return nil })

		res, err := rtt.RoundTrip(httptest.NewRequest(http.MethodGet, srv.URL, nil))
		require.NoError(t, err)

		select {
		case <-retryCh:
			t.Fatal("retry called, but shouldn't ")
		case <-time.After(time.Second):
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)

	})
}
