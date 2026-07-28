package checkpolicies

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/policystate"
)

func TestUseCase_Execute_PrintsPolicyState(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanPolicyState", ctx, projectID, scanID).Return(policystate.Rejected, nil).Once()

	cli := NewMockCLI(t)
	cli.On("ReturnText", ctx, "Rejected").Return().Once()

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_FailOnPoliciesRejected(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanPolicyState", ctx, projectID, scanID).Return(policystate.Rejected, nil).Once()

	cli := NewMockCLI(t)
	cli.On("ReturnText", ctx, "Rejected").Return().Once()

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)

	err = uc.Execute(ctx, scanID, true)
	require.Error(t, err)

	var failErr *apperror.FailError
	require.ErrorAs(t, err, &failErr)
}

func TestUseCase_Execute_FailOnPoliciesRejected_NoneAndConfirmedOk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state policystate.State
	}{
		{name: "confirmed", state: policystate.Confirmed},
		{name: "none", state: policystate.None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			projectID := uuid.New()
			scanID := uuid.New()
			cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

			ai := NewMockAI(t)
			ai.On("InitializeWithRetry", ctx).Return(nil).Once()
			ai.On("GetScanPolicyState", ctx, projectID, scanID).Return(tt.state, nil).Once()

			cli := NewMockCLI(t)
			cli.On("ReturnText", ctx, tt.state.String()).Return().Once()

			uc, err := NewUseCase(ai, cli, cfg)
			require.NoError(t, err)
			require.NoError(t, uc.Execute(ctx, scanID, true))
		})
	}
}
