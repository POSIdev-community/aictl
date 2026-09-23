package sbomproject_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/POSIdev-community/aictl/internal/core/usecase/get/scans/sbomproject"
)

type stubAI struct {
	version     version.Version
	project     *project.Project
	scans       []scan.Scan
	gotBranchID uuid.UUID
}

func (s *stubAI) InitializeWithRetry(context.Context) error { return nil }
func (s *stubAI) GetVersion(context.Context) (version.Version, error) {
	return s.version, nil
}
func (s *stubAI) GetProject(context.Context, uuid.UUID) (*project.Project, error) {
	return s.project, nil
}
func (s *stubAI) GetScans(_ context.Context, branchId uuid.UUID) ([]scan.Scan, error) {
	s.gotBranchID = branchId
	return s.scans, nil
}
func (s *stubAI) GetLastScan(_ context.Context, branchId uuid.UUID) (*scan.Scan, error) {
	s.gotBranchID = branchId
	return &s.scans[0], nil
}

type stubCLI struct{ shown []scan.Scan }

func (s *stubCLI) ShowScans(_ context.Context, items []scan.Scan)      { s.shown = items }
func (s *stubCLI) ShowScansQuite(_ context.Context, items []scan.Scan) { s.shown = items }

func emptyFilter(t *testing.T) regexfilter.RegexFilter {
	t.Helper()
	f, err := regexfilter.NewRegexFilter("")
	require.NoError(t, err)
	return f
}

func mustVersion(t *testing.T, raw string) version.Version {
	t.Helper()
	v, err := version.NewVersion(raw)
	require.NoError(t, err)
	return v
}

func TestExecute_ResolvesVirtualBranch(t *testing.T) {
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	virtualBranchID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	scanID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	ai := &stubAI{
		version: mustVersion(t, "6.3.0"),
		project: &project.Project{
			Id: projectID, Name: "sbom", Type: project.TypeSbom, VirtualBranchId: virtualBranchID,
		},
		scans: []scan.Scan{scan.NewScan(scanID, uuid.Nil, nil, nil)},
	}
	cli := &stubCLI{}
	uc, err := sbomproject.NewUseCase(ai, cli, config.NewConfig(config.Uri{}, "", true, projectID, uuid.Nil))
	require.NoError(t, err)

	require.NoError(t, uc.Execute(context.Background(), emptyFilter(t), false, false))
	require.Equal(t, virtualBranchID, ai.gotBranchID)
	require.Len(t, cli.shown, 1)
}

func TestExecute_RejectsSource(t *testing.T) {
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ai := &stubAI{
		version: mustVersion(t, "6.3.0"),
		project: &project.Project{Id: projectID, Type: project.TypeSource},
	}
	uc, err := sbomproject.NewUseCase(ai, &stubCLI{}, config.NewConfig(config.Uri{}, "", true, projectID, uuid.Nil))
	require.NoError(t, err)

	err = uc.Execute(context.Background(), emptyFilter(t), false, false)
	require.ErrorIs(t, err, project.ErrCannotGetSbomScansOnSource)
}

func TestExecute_UnsupportedBefore63(t *testing.T) {
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ai := &stubAI{version: mustVersion(t, "6.2.0")}
	uc, err := sbomproject.NewUseCase(ai, &stubCLI{}, config.NewConfig(config.Uri{}, "", true, projectID, uuid.Nil))
	require.NoError(t, err)

	err = uc.Execute(context.Background(), emptyFilter(t), false, false)
	require.ErrorContains(t, err, project.ErrSbomUnsupported)
}
