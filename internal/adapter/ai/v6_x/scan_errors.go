package v6_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
)

func (a *ClientAI61) GetScanErrors(ctx context.Context, projectId, scanResultId uuid.UUID) ([]string, error) {
	response, err := a.GetApiProjectsProjectIdScanResultsScanResultIdErrorsWithResponse(
		ctx, projectId, scanResultId, a.AddJWTToHeader,
	)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan errors request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, response.JSON400); err != nil {
		return nil, fmt.Errorf("ai adapter get scan errors: %w", err)
	}

	if response.JSON200 == nil {
		return []string{}, nil
	}

	lines := make([]string, 0, len(*response.JSON200))
	for _, model := range *response.JSON200 {
		line := common.GetOrDefault(model.Message, "")
		if line == "" {
			line = common.GetOrDefault(model.ErrorType, "")
		}
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
