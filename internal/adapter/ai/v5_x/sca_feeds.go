package v5_x

import (
	"context"
	"fmt"
	"io"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func (a *ClientAI5x) UpdateScaFeeds(_ context.Context, _, _ string) error {
	return fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}

func (a *ClientAI5x) GetScaFeeds(_ context.Context, _ []scafeeds.Status) ([]scafeeds.Package, error) {
	return nil, fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}

func (a *ClientAI5x) DownloadScaFeeds(_ context.Context, _ string) (io.ReadCloser, string, error) {
	return nil, "", fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}

func (a *ClientAI5x) RollbackScaFeeds(_ context.Context) (scafeeds.Package, error) {
	return scafeeds.Package{}, fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}
