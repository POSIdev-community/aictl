package ai

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common/notify"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestScanStageFromNotificationProgress(t *testing.T) {
	t.Parallel()

	scanID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	stageName := "Scan"
	value := int32(10)
	msg := notify.Message{
		ScanProgress: &notify.ScanProgressNotification{
			ScanResultID: scanID,
			Progress: &notify.ProgressModel{
				Stage: &stageName,
				Value: &value,
			},
		},
	}

	stage, ok := scanStageFromNotification(msg, scanID)
	require.True(t, ok)
	require.Equal(t, "Scan", stage.Stage)
	require.Equal(t, int32(10), stage.Value)

	_, ok = scanStageFromNotification(msg, uuid.New())
	require.False(t, ok)
}

func TestScanStageFromNotificationCompleted(t *testing.T) {
	t.Parallel()

	scanID := uuid.New()
	done := "Failed"
	msg := notify.Message{
		ScanCompleted: &notify.ScanCompletedNotification{
			ScanResultID: scanID,
			Stage:        &done,
		},
	}

	stage, ok := scanStageFromNotification(msg, scanID)
	require.True(t, ok)
	require.Equal(t, "Failed", stage.Stage)
}
