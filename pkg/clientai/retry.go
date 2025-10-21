package clientai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type retryHandler func(context.Context, *http.Request) error

func NewRetryRoundTripper(rt http.RoundTripper, onCode int, method retryHandler) *RetryRoundTripper {
	return &RetryRoundTripper{rt: rt, onCode: onCode, method: method}
}

// RetryRoundTripper is a custom HTTP transport that retries requests based on specific status codes and retry handlers.
type RetryRoundTripper struct {
	rt     http.RoundTripper
	onCode int
	method retryHandler
}

// RoundTrip executes a single HTTP transaction, retry if the response status code matches the specified code, invoking a retry handler before retry.
func (rrt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rrt.rt == nil {
		rrt.rt = http.DefaultTransport
	}

	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}

		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	resp, err := rrt.rt.RoundTrip(req)
	if err != nil || resp.StatusCode != rrt.onCode {
		return resp, err
	}

	_ = resp.Body.Close()

	clonedReq := req.Clone(req.Context())
	clonedReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if rrt.method != nil {
		if err := rrt.method(req.Context(), clonedReq); err != nil {
			return resp, fmt.Errorf("call retry handler: %w", err)
		}
	}

	return rrt.rt.RoundTrip(clonedReq)
}
