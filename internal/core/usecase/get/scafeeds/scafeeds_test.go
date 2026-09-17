package scafeeds_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	getscafeeds "github.com/POSIdev-community/aictl/internal/core/usecase/get/scafeeds"
)

func TestGetScaFeedsUseCase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		ai := getscafeeds.NewMockAI(t)
		cli := getscafeeds.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		pkgs := []scafeeds.Package{{Version: "1", Status: scafeeds.StatusCurrent}}
		ai.On("GetScaFeeds", ctx, []scafeeds.Status{scafeeds.StatusCurrent}).Return(pkgs, nil).Once()
		cli.On("ShowScaFeeds", ctx, pkgs).Return().Once()

		uc, err := getscafeeds.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, []scafeeds.Status{scafeeds.StatusCurrent}))
	})

	t.Run("get_error", func(t *testing.T) {
		t.Parallel()
		ai := getscafeeds.NewMockAI(t)
		cli := getscafeeds.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("GetScaFeeds", ctx, mock.Anything).Return(nil, fmt.Errorf("boom")).Once()

		uc, err := getscafeeds.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.Error(t, uc.Execute(ctx, nil))
	})
}
