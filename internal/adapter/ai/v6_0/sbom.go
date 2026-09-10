package v6_0

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
)

func (a *ClientAI60) CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	return nil, fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI60) UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error {
	return fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI60) StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	return uuid.UUID{}, fmt.Errorf("%s", project.ErrSbomUnsupported)
}

func (a *ClientAI60) GetProjectByName(ctx context.Context, projectName string) (*project.Project, error) {
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
