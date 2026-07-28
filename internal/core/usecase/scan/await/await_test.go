package await

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

const testPollInterval = 20 * time.Millisecond

func TestAwaitQueueDisplay(t *testing.T) {
	t.Parallel()

	itemWithPlace := queue.Item{Place: 2, OutOf: 5, ScanId: uuid.New()}
	if itemWithPlace.OutOf <= 0 {
		t.Fatal("expected queue item with position")
	}

	itemWithoutPlace := queue.Item{ScanId: uuid.New()}
	if itemWithoutPlace.OutOf != 0 {
		t.Fatal("expected queue item without position when scan left queue")
	}
}

func TestUseCase_Execute_SilenceFallbackDone(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	updates := make(chan scanstage.ScanStage)

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: "Scanning", Value: 10}, nil,
	).Once()
	ai.On("WatchScanStage", ctx, scanID).Return((<-chan scanstage.ScanStage)(updates), nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: scanstage.Done}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Done).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, testPollInterval)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_SilenceFallbackStillRunningThenDone(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	updates := make(chan scanstage.ScanStage)

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: "Scanning", Value: 10}, nil,
	).Once()
	ai.On("WatchScanStage", ctx, scanID).Return((<-chan scanstage.ScanStage)(updates), nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: "Scanning", Value: 50}, nil,
	).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: scanstage.Done}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Done).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, testPollInterval)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_NotificationCompletesBeforeSilence(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	updates := make(chan scanstage.ScanStage, 1)
	updates <- scanstage.ScanStage{Stage: scanstage.Done}

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: "Scanning", Value: 10}, nil,
	).Once()
	ai.On("WatchScanStage", ctx, scanID).Return((<-chan scanstage.ScanStage)(updates), nil).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Done).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, time.Hour)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_SilenceFallbackErrorThenDone(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	updates := make(chan scanstage.ScanStage)

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: "Scanning", Value: 10}, nil,
	).Once()
	ai.On("WatchScanStage", ctx, scanID).Return((<-chan scanstage.ScanStage)(updates), nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{}, errors.New("temporary"),
	).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: scanstage.Done}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Done).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, testPollInterval)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_AlreadyComplete(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: scanstage.Done}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Done).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, testPollInterval)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, scanID, false))
}

func TestUseCase_Execute_FailOnScanFailed(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	projectID := uuid.New()
	scanID := uuid.New()
	cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetScanStage", ctx, projectID, scanID).Return(
		scanstage.ScanStage{Stage: scanstage.Failed}, nil,
	).Once()

	cli := NewMockCLI(t)
	cli.On("ShowTextf", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ShowTextf", mock.Anything, mock.Anything).Return().Maybe()
	cli.On("ReturnText", ctx, scanstage.Failed).Return().Once()

	uc, err := NewUseCase(ai, cli, cfg, testPollInterval)
	require.NoError(t, err)

	err = uc.Execute(ctx, scanID, true)
	require.Error(t, err)

	var failErr *apperror.FailError
	require.ErrorAs(t, err, &failErr)
}

func TestNewUseCase_InvalidPollInterval(t *testing.T) {
	t.Parallel()

	_, err := NewUseCase(NewMockAI(t), NewMockCLI(t), nil, 0)
	require.Error(t, err)
}
