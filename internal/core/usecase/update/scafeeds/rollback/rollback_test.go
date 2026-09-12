package rollback_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/usecase/update/scafeeds/rollback"
)

func TestRollbackScaFeedsUseCase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("yes_skip_confirm", func(t *testing.T) {
		t.Parallel()
		ai := rollback.NewMockAI(t)
		cli := rollback.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		ai.On("RollbackScaFeeds", ctx).Return(scafeeds.Package{Version: "1.2.0"}, nil).Once()
		cli.On("ReturnText", ctx, "1.2.0").Return().Once()
		cli.On("ShowTextf", ctx, "SCA feeds rolled back to version '%s'", mock.Anything).Return().Once()

		uc, err := rollback.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, true))
	})

	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()
		ai := rollback.NewMockAI(t)
		cli := rollback.NewMockCLI(t)
		ai.On("InitializeWithRetry", ctx).Return(nil).Once()
		cli.On("AskConfirmation", ctx, "Roll back current SCA feeds package?").Return(false, nil).Once()
		cli.On("ShowText", ctx, "Cancelled").Return().Once()

		uc, err := rollback.NewUseCase(ai, cli)
		require.NoError(t, err)
		require.NoError(t, uc.Execute(ctx, false))
	})
}
