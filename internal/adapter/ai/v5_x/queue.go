package v5_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/pkg/clientai/v5_x"
)

func (a *ClientAI5x) GetScanQueue(ctx context.Context) ([]queue.Entry, error) {
	response, err := a.GetApiScansWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan queue request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get scan queue: %w", err)
	}

	if response.JSON200 == nil {
		return []queue.Entry{}, nil
	}

	entries := make([]queue.Entry, 0, len(*response.JSON200))
	for _, model := range *response.JSON200 {
		entry := queue.Entry{Stage: scanstage.Enqueued}
		if model.ProjectId != nil {
			entry.ProjectId = uuid.UUID(*model.ProjectId)
		}
		if model.BranchId != nil {
			entry.BranchId = uuid.UUID(*model.BranchId)
		}
		if model.ScanResultId != nil {
			entry.ScanId = uuid.UUID(*model.ScanResultId)
		}
		if model.StatusType != nil {
			entry.Stage = mapScanStatusType(*model.StatusType)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (a *ClientAI5x) GetActiveScans(ctx context.Context) ([]queue.Entry, error) {
	response, err := a.GetApiProjectsActiveScansWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get active scans request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get active scans: %w", err)
	}

	if response.JSON200 == nil {
		return []queue.Entry{}, nil
	}

	entries := make([]queue.Entry, 0, len(*response.JSON200))
	for _, model := range *response.JSON200 {
		entry := queue.Entry{}
		if model.Project != nil && model.Project.Id != nil {
			entry.ProjectId = uuid.UUID(*model.Project.Id)
		}
		if model.Branch != nil && model.Branch.Id != nil {
			entry.BranchId = uuid.UUID(*model.Branch.Id)
		}
		if model.ScanResultId != nil {
			entry.ScanId = uuid.UUID(*model.ScanResultId)
		}
		if model.Progress != nil && model.Progress.Stage != nil {
			entry.Stage = string(*model.Progress.Stage)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// cancelActiveScan stops a scan that already left the queue and is running on an agent.
// Returns cancelled=true when an agent pause/stop was issued.
func (a *ClientAI5x) cancelActiveScan(ctx context.Context, scanResultId uuid.UUID) (bool, error) {
	response, err := a.GetApiScanAgentsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai adapter get scan agents request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return false, fmt.Errorf("ai adapter get scan agents: %w", err)
	}
	if response.JSON200 == nil {
		return false, nil
	}

	for _, agent := range *response.JSON200 {
		if agent.ScanResultId == nil || uuid.UUID(*agent.ScanResultId) != scanResultId {
			continue
		}
		if agent.Id == nil {
			continue
		}

		if err := a.cancelScanOnAgent(ctx, uuid.UUID(*agent.Id)); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (a *ClientAI5x) cancelScanOnAgent(ctx context.Context, agentId uuid.UUID) error {
	stopScan := true
	params := &v5_x.PostApiScanAgentsScanAgentIdPauseParams{StopScan: &stopScan}
	response, err := a.PostApiScanAgentsScanAgentIdPauseWithResponse(ctx, agentId, params, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter cancel scan on agent request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return fmt.Errorf("ai adapter cancel scan on agent: %w", err)
	}

	return nil
}

func mapScanStatusType(status v5_x.ScanStatusType) string {
	switch status {
	case v5_x.ScanStatusTypePending, v5_x.ScanStatusTypeScheduled:
		return scanstage.Enqueued
	case v5_x.ScanStatusTypeScan:
		return "Scan"
	case v5_x.ScanStatusTypeFinished:
		return scanstage.Done
	case v5_x.ScanStatusTypeFailed:
		return scanstage.Failed
	case v5_x.ScanStatusTypeAborted:
		return scanstage.Aborted
	default:
		return string(status)
	}
}
