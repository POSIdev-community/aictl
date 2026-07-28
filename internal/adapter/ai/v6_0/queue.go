package v6_0

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
)

func (a *ClientAI60) GetScanQueue(ctx context.Context) ([]queue.Entry, error) {
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

func (a *ClientAI60) GetActiveScans(ctx context.Context) ([]queue.Entry, error) {
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
