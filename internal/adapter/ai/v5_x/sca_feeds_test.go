package v5_x_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func TestUpdateScaFeedsUnsupported(t *testing.T) {
	t.Parallel()

	client := &v5_x.ClientAI5x{}
	err := client.UpdateScaFeeds(t.Context(), "./feeds.zip", "1")
	require.Error(t, err)
	require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
}
