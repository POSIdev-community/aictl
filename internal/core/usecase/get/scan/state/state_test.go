package state

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

func TestUseCase_Execute_FailOnScanFailed(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: Failed}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ReturnText", ctx, Failed).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)

	err = uc.Execute(ctx, scanID, true)
	require.Error(t, err)

	var failErr *apperror.FailError
	require.ErrorAs(t, err, &failErr)
}

func TestUseCase_Execute_WithoutFailFlag(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: Failed}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ReturnText", ctx, Failed).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}
