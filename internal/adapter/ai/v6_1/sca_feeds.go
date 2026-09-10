package v6_1

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func (a *ClientAI61) UpdateScaFeeds(_ context.Context, _, _ string) error {
	return fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)
}
