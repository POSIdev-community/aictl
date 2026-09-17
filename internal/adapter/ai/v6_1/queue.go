package v6_1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

func (a *ClientAI61) GetScanQueue(ctx context.Context) ([]queue.Entry, error) {
	response, err := a.GetAllItemsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan queue request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return nil, fmt.Errorf("ai adapter get scan queue: %w", err)
	}

	if response.JSON200 == nil {
		return []queue.Entry{}, nil
	}

	entries := make([]queue.Entry, 0, len(*response.JSON200))
	for _, model := range *response.JSON200 {
		entries = append(entries, queue.Entry{
			ProjectId: model.ScanObject.ProjectId,
			BranchId:  model.ScanObject.BranchId,
			ScanId:    model.ScanResultId,
			Stage:     scanstage.Enqueued,
		})
	}

	return entries, nil
}

func (a *ClientAI61) GetActiveScans(ctx context.Context) ([]queue.Entry, error) {
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
// Returns cancelled=true when an agent cancel was issued.
func (a *ClientAI61) cancelActiveScan(ctx context.Context, scanResultId uuid.UUID) (bool, error) {
	active, err := a.GetActiveScans(ctx)
	if err != nil {
		return false, fmt.Errorf("ai adapter cancel active scan: %w", err)
	}

	var projectId, branchId uuid.UUID
	found := false
	for _, e := range active {
		if e.ScanId == scanResultId {
			projectId = e.ProjectId
			branchId = e.BranchId
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	response, err := a.GetAllWithScanInfoWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return false, fmt.Errorf("ai adapter get agents with scans request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return false, fmt.Errorf("ai adapter get agents with scans: %w", err)
	}
	if response.JSON200 == nil {
		return false, nil
	}

	for _, item := range *response.JSON200 {
		if item.Scan == nil {
			continue
		}
		if item.Scan.Object.ProjectId == projectId && item.Scan.Object.BranchId == branchId {
			if err := a.cancelScanOnAgent(ctx, item.Agent.Id); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	// Active entry exists but agent already released it.
	return true, nil
}

func (a *ClientAI61) cancelScanOnAgent(ctx context.Context, agentId uuid.UUID) error {
	response, err := a.CancelScanWithResponse(ctx, agentId, a.AddJWTToHeader)
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
