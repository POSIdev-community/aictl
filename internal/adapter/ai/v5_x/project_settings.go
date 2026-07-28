package v5_x

import (
	"context"
	"fmt"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/client"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/google/uuid"
)

func (a *ClientAI5x) GetProjectSettings(ctx context.Context, projectId uuid.UUID) (settings.ScanSettings, error) {
	_ = ctx
	_ = projectId

	return settings.ScanSettings{}, client.ErrProjectScanSettingsNotSupported()
}

func (a *ClientAI5x) GetScanAgents(ctx context.Context) ([]scanagent.ScanAgent, error) {
	response, err := a.GetApiScanAgentsWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	if err = CheckResponseByModel(statusCode, body, nil); err != nil {
		return nil, fmt.Errorf("ai adapter get scan agents: %w", err)
	}

	if response.JSON200 == nil {
		return []scanagent.ScanAgent{}, nil
	}

	models := *response.JSON200
	agents := make([]scanagent.ScanAgent, 0, len(models))
	for _, model := range models {
		if model.Id == nil {
			continue
		}

		status := ""
		if model.StatusType != nil {
			status = string(*model.StatusType)
		}

		agents = append(agents, scanagent.ScanAgent{
			Id:              uuid.UUID(*model.Id),
			Name:            common.GetOrDefault(model.Name, ""),
			Status:          status,
			Version:         common.GetOrDefault(model.Version, ""),
			OperatingSystem: common.GetOrDefault(model.OperatingSystem, ""),
		})
	}

	return agents, nil
}
