package v6_1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_1"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func TestUpdateScaFeedsUnsupported(t *testing.T) {
	t.Parallel()

	client := &v6_1.ClientAI61{}
	err := client.UpdateScaFeeds(t.Context(), "./feeds.zip", "1")
	require.Error(t, err)
	require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
}
