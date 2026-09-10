package v5_x

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func (a *ClientAI5x) UpdateScaFeeds(_ context.Context, _, _ string) error {
	return fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}
