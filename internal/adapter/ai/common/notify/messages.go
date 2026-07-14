package notify

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// SignalR message types (JSON protocol).
const (
	messageTypeInvocation = 1
	messageTypePing       = 6
	messageTypeClose      = 7
)

// Hub notification targets used for scan await.
const (
	TargetScanProgress        = "ScanProgress"
	TargetScanCompleted       = "ScanCompleted"
	TargetNeedRefreshToken    = "NeedRefreshToken"
	TargetNeedSyncClientState = "NeedSyncClientState"
)

// Message is a decoded SignalR hub frame.
type Message struct {
	Type                int
	Target              string
	Arguments           []json.RawMessage
	Error               string
	NeedRefreshToken    bool
	NeedSyncClientState bool
	ScanProgress        *ScanProgressNotification
	ScanCompleted       *ScanCompletedNotification
}

type ProgressModel struct {
	Stage    *string `json:"stage"`
	SubStage *string `json:"subStage"`
	Value    *int32  `json:"value"`
}

type ScanProgressNotification struct {
	ProjectID    uuid.UUID      `json:"projectId"`
	BranchID     uuid.UUID      `json:"branchId"`
	ScanResultID uuid.UUID      `json:"scanResultId"`
	Progress     *ProgressModel `json:"progress"`
}

type ScanCompletedNotification struct {
	ProjectID    uuid.UUID `json:"projectId"`
	BranchID     uuid.UUID `json:"branchId"`
	ScanResultID uuid.UUID `json:"scanResultId"`
	Stage        *string   `json:"stage"`
}

func parseFrame(frame []byte) (Message, bool, error) {
	var raw struct {
		Type      int               `json:"type"`
		Target    string            `json:"target"`
		Arguments []json.RawMessage `json:"arguments"`
		Error     string            `json:"error"`
	}
	if err := json.Unmarshal(frame, &raw); err != nil {
		return Message{}, false, fmt.Errorf("decode frame: %w", err)
	}

	msg := Message{
		Type:      raw.Type,
		Target:    raw.Target,
		Arguments: raw.Arguments,
		Error:     raw.Error,
	}

	if raw.Type != messageTypeInvocation {
		return msg, true, nil
	}

	if raw.Target == TargetNeedRefreshToken {
		msg.NeedRefreshToken = true

		return msg, true, nil
	}

	if raw.Target == TargetNeedSyncClientState {
		msg.NeedSyncClientState = true

		return msg, true, nil
	}

	if len(raw.Arguments) == 0 {
		return msg, true, nil
	}

	switch raw.Target {
	case TargetScanProgress:
		var n ScanProgressNotification
		if err := json.Unmarshal(raw.Arguments[0], &n); err != nil {
			return msg, true, nil
		}
		msg.ScanProgress = &n
	case TargetScanCompleted:
		var n ScanCompletedNotification
		if err := json.Unmarshal(raw.Arguments[0], &n); err != nil {
			return msg, true, nil
		}
		msg.ScanCompleted = &n
	}

	return msg, true, nil
}
