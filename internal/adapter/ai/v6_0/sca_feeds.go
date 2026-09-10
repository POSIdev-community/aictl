package v6_0

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func (a *ClientAI60) UpdateScaFeeds(_ context.Context, _, _ string) error {
	return fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}
