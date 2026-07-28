package v6_0

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_0"
)

func (a *ClientAI60) GetReportTemplates(ctx context.Context, localization string) ([]report.Template, error) {
	withContent := false
	params := &v6_0.GetApiReportsTemplatesParams{
		LocaleId:    &localization,
		WithContent: &withContent,
	}

	response, err := a.GetApiReportsTemplates(ctx, params, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get report templates request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if err = CheckResponse(response, "report templates"); err != nil {
		return nil, fmt.Errorf("ai adapter get report templates: %w", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ai adapter read report templates body: %w", err)
	}

	// Minimal model: full ReportTemplateModel fails on zero creationDate without timezone.
	type templateDTO struct {
		Id   *uuid.UUID `json:"id,omitempty"`
		Name *string    `json:"name,omitempty"`
	}

	var models []templateDTO
	if err := json.Unmarshal(bodyBytes, &models); err != nil {
		return nil, fmt.Errorf("ai adapter decode report templates: %w", err)
	}

	templates := make([]report.Template, 0, len(models))
	for _, model := range models {
		if model.Id == nil {
			continue
		}
		templates = append(templates, report.Template{
			Id:   *model.Id,
			Name: common.GetOrDefault(model.Name, ""),
		})
	}

	return templates, nil
}
