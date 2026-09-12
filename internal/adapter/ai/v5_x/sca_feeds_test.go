package v5_x_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func TestScaFeedsUnsupported(t *testing.T) {
	t.Parallel()

	client := &v5_x.ClientAI5x{}

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		err := client.UpdateScaFeeds(t.Context(), "./feeds.zip", "1")
		require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
	})

	t.Run("get", func(t *testing.T) {
		t.Parallel()
		_, err := client.GetScaFeeds(t.Context(), nil)
		require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
	})

	t.Run("download", func(t *testing.T) {
		t.Parallel()
		_, _, err := client.DownloadScaFeeds(t.Context(), "1")
		require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
	})

	t.Run("rollback", func(t *testing.T) {
		t.Parallel()
		_, err := client.RollbackScaFeeds(t.Context())
		require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
	})
}
