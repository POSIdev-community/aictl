package v5_x

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
)

func (a *ClientAI5x) CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	return nil, fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI5x) UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error {
	return fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI5x) StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	return uuid.UUID{}, fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI5x) GetProjectByName(ctx context.Context, projectName string) (*project.Project, error) {
	id, err := a.GetProjectId(ctx, projectName)
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, nil
	}
	p := project.NewProjectWithType(*id, projectName, project.TypeSource)

	return &p, nil
}
