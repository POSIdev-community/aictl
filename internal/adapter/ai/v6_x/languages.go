package v6_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (a *ClientAI61) UpdateProjectLanguages(ctx context.Context, projectId uuid.UUID) error {
	defaultRes, err := a.GetApiProjectsProjectIdDefaultSettingsWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("ai adapter get project default settings request: %w", err)
	}

	statusCode := defaultRes.StatusCode()
	body := string(defaultRes.Body)
	if err = CheckResponseByModel(statusCode, body, defaultRes.JSON400); err != nil {
		return fmt.Errorf("ai adapter get project default settings: %w", err)
	}

	if defaultRes.JSON200 == nil {
		return fmt.Errorf("ai adapter get project default settings: empty response")
	}

	languages := []string{}
	if defaultRes.JSON200.Languages != nil {
		languages = make([]string, len(*defaultRes.JSON200.Languages))
		for i, lang := range *defaultRes.JSON200.Languages {
			languages[i] = string(lang)
		}
	}

	current, err := a.GetProjectSettings(ctx, projectId)
	if err != nil {
		return fmt.Errorf("get project settings: %w", err)
	}

	current.Languages = languages

	if err := a.SetProjectSettings(ctx, projectId, &current); err != nil {
		return fmt.Errorf("set project settings: %w", err)
	}

	return nil
}
