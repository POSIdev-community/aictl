package common

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

type RetryRoundTripper struct {
	rt     http.RoundTripper
	onCode int
	method retryHandler
}

func (rrt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rrt.rt == nil {
		rrt.rt = http.DefaultTransport
	}

	var bodyBytes []byte
	if req.Body != nil && req.Body != http.NoBody {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	resp, err := rrt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != rrt.onCode {
		return resp, nil
	}

	_ = resp.Body.Close()

	clonedReq := req.Clone(req.Context())
	if bodyBytes != nil {
		clonedReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	} else {
		clonedReq.Body = nil
	}

	if rrt.method != nil {
		if err := rrt.method(req.Context(), clonedReq); err != nil {
			return nil, fmt.Errorf("call retry handler: %w", err)
		}
	}

	return rrt.rt.RoundTrip(clonedReq)
}
