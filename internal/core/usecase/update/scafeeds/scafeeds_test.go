package scafeeds_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	usecase "github.com/POSIdev-community/aictl/internal/core/usecase/update/scafeeds"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		ai := usecase.NewMockAI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("UpdateScaFeeds", ctx, "./feeds.zip", "47").Return(nil).Once()

		cli := usecase.NewMockCLI(t)
		cli.On("ShowTextf", ctx, "SCA feeds updated, version '%s'", mock.Anything).Return().Once()

		uc, err := usecase.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, "./feeds.zip", "47"))
	})

	t.Run("unsupported_version", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		ai := usecase.NewMockAI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("UpdateScaFeeds", ctx, "./feeds.zip", "1").Return(fmt.Errorf("%s", scafeeds.ErrScaFeedsUnsupported)).Once()

		cli := usecase.NewMockCLI(t)

		uc, err := usecase.NewUseCase(ai, cli)
		require.NoError(t, err)
		err = uc.Execute(ctx, "./feeds.zip", "1")
		require.Error(t, err)
		require.ErrorContains(t, err, scafeeds.ErrScaFeedsUnsupported)
	})

	t.Run("nil_adapters", func(t *testing.T) {
		t.Parallel()

		_, err := usecase.NewUseCase(nil, usecase.NewMockCLI(t))
		require.Error(t, err)

		_, err = usecase.NewUseCase(usecase.NewMockAI(t), nil)
		require.Error(t, err)
	})
}
