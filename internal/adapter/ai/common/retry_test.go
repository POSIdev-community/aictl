package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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
		t.Cleanup(srv.Close)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error {
			close(retryCh)

			return nil
		})

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
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
			}),
		)
		t.Cleanup(srv.Close)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error {
			close(retryCh)

			return http.ErrHandlerTimeout
		})

		res, err := rtt.RoundTrip(httptest.NewRequest(http.MethodGet, srv.URL, nil))
		require.Error(t, err)
		assert.ErrorIs(t, err, http.ErrHandlerTimeout)
		assert.Nil(t, res)

		select {
		case <-retryCh:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for retry")
		}
	})

	t.Run("should not retry on success", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)
		t.Cleanup(srv.Close)

		retryCh := make(chan struct{})

		rtt := NewRetryRoundTripper(&http.Transport{}, http.StatusUnauthorized, func(ctx context.Context, req *http.Request) error {
			close(retryCh)

			return nil
		})

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

func TestDoJWTRefresh_WaitsForInFlightRefresh(t *testing.T) {
	t.Parallel()

	base := NewBaseClient()

	var calls atomic.Int32
	var wg sync.WaitGroup
	const n = 16

	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			err := base.DoJWTRefresh(func() error {
				calls.Add(1)
				time.Sleep(50 * time.Millisecond)

				return nil
			})
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), calls.Load())
}
