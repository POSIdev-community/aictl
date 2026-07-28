package v5_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/pkg/clientai/v5_x"
)

func (a *ClientAI5x) GetProjectExclusions(ctx context.Context, projectId uuid.UUID) (string, error) {
	response, err := a.GetApiProjectsProjectIdExclusionsWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return "", fmt.Errorf("ai adapter get project exclusions request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return "", fmt.Errorf("ai adapter get project exclusions: %w", err)
	}

	if response.JSON200 == nil || response.JSON200.Exclusions == nil {
		return "", nil
	}

	return *response.JSON200.Exclusions, nil
}

func (a *ClientAI5x) SetProjectExclusions(ctx context.Context, projectId uuid.UUID, exclusions string) error {
	body := v5_x.FileFolderExclusionsModel{Exclusions: &exclusions}
	response, err := a.PutApiProjectsProjectIdExclusionsWithResponse(ctx, projectId, body, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter set project exclusions request: %w", err)
	}

	statusCode := response.StatusCode()
	respBody := string(response.Body)
	if err = CheckResponseByModel(statusCode, respBody, response.JSON400); err != nil {
		return fmt.Errorf("ai adapter set project exclusions: %w", err)
	}

	return nil
}
