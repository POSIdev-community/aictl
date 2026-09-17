package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
)

type ProjectGetter interface {
	GetProject(ctx context.Context, projectId uuid.UUID) (*project.Project, error)
}

func RequireSourceProject(ctx context.Context, ai ProjectGetter, projectId uuid.UUID, sbomErr error) error {
	p, err := ai.GetProject(ctx, projectId)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}
	if p.Type.IsSbom() {
		return sbomErr
	}

	return nil
}

func RequireSbomProject(ctx context.Context, ai ProjectGetter, projectId uuid.UUID, sourceErr error) error {
	p, err := ai.GetProject(ctx, projectId)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}
	if !p.Type.IsSbom() {
		return sourceErr
	}

	return nil
}
