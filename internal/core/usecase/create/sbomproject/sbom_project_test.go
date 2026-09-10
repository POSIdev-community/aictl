package sbomproject_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/usecase/create/sbomproject"
)

type stubAI struct {
	byName    *project.Project
	createdID *uuid.UUID
	updated   bool
}

func (s *stubAI) InitializeWithRetry(ctx context.Context) error { return nil }
func (s *stubAI) GetProjectByName(ctx context.Context, projectName string) (*project.Project, error) {
	return s.byName, nil
}
func (s *stubAI) CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	s.createdID = &id
	return s.createdID, nil
}
func (s *stubAI) UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error {
	s.updated = true
	return nil
}

type stubCLI struct {
	out string
}

func (s *stubCLI) ReturnText(ctx context.Context, text string)            { s.out = text }
func (s *stubCLI) ShowTextf(ctx context.Context, format string, a ...any) {}

func TestCreateSbomProject_SafeReuploads(t *testing.T) {
	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	ai := &stubAI{byName: &project.Project{Id: id, Name: "x", Type: project.TypeSbom}}
	cli := &stubCLI{}
	uc, err := sbomproject.NewUseCase(ai, cli)
	require.NoError(t, err)

	require.NoError(t, uc.Execute(context.Background(), "x", "/tmp/sbom.json", true))
	require.True(t, ai.updated)
	require.Equal(t, id.String(), cli.out)
}

func TestCreateSbomProject_RejectsSourceName(t *testing.T) {
	id := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	ai := &stubAI{byName: &project.Project{Id: id, Name: "x", Type: project.TypeSource}}
	uc, err := sbomproject.NewUseCase(ai, &stubCLI{})
	require.NoError(t, err)

	err = uc.Execute(context.Background(), "x", "/tmp/sbom.json", true)
	require.ErrorIs(t, err, project.ErrNameUsedBySource)
}
