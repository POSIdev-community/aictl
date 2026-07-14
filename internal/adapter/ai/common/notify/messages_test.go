package notify

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseFrameNeedSyncClientState(t *testing.T) {
	t.Parallel()

	frame := []byte(`{"type":1,"target":"NeedSyncClientState","arguments":[]}`)

	msg, ok, err := parseFrame(frame)
	require.NoError(t, err)
	require.True(t, ok)
	require.True(t, msg.NeedSyncClientState)
	require.Equal(t, TargetNeedSyncClientState, msg.Target)
}

func TestParseFrameNeedRefreshToken(t *testing.T) {
	t.Parallel()

	frame := []byte(`{"type":1,"target":"NeedRefreshToken","arguments":[]}`)

	msg, ok, err := parseFrame(frame)
	require.NoError(t, err)
	require.True(t, ok)
	require.True(t, msg.NeedRefreshToken)
	require.Equal(t, TargetNeedRefreshToken, msg.Target)
}

func TestParseFrameScanProgress(t *testing.T) {
	t.Parallel()

	scanID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	frame := []byte(`{"type":1,"target":"ScanProgress","arguments":[{"scanResultId":"11111111-1111-1111-1111-111111111111","progress":{"stage":"Scan","value":42,"subStage":"engine"}}]}`)

	msg, ok, err := parseFrame(frame)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, TargetScanProgress, msg.Target)
	require.NotNil(t, msg.ScanProgress)
	require.Equal(t, scanID, msg.ScanProgress.ScanResultID)
	require.NotNil(t, msg.ScanProgress.Progress)
	require.Equal(t, "Scan", *msg.ScanProgress.Progress.Stage)
	require.Equal(t, int32(42), *msg.ScanProgress.Progress.Value)
	require.Equal(t, "engine", *msg.ScanProgress.Progress.SubStage)
}

func TestParseFrameScanCompleted(t *testing.T) {
	t.Parallel()

	frame := []byte(`{"type":1,"target":"ScanCompleted","arguments":[{"scanResultId":"11111111-1111-1111-1111-111111111111","stage":"Done"}]}`)

	msg, ok, err := parseFrame(frame)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotNil(t, msg.ScanCompleted)
	require.Equal(t, "Done", *msg.ScanCompleted.Stage)
}

func TestSplitRecords(t *testing.T) {
	t.Parallel()

	parts := splitRecords([]byte("a\x1eb\x1e"))
	require.Equal(t, [][]byte{[]byte("a"), []byte("b")}, parts)
}
