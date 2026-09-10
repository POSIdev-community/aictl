package v6_0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func TestUpdateScaFeedsUnsupported(t *testing.T) {
	t.Parallel()

	client := &v6_0.ClientAI60{}
	err := client.UpdateScaFeeds(t.Context(), "./feeds.zip", "1")
	require.Error(t, err)
	require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
}
