package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeProjectPolicyCheckUC struct{ called int }

func (f *fakeProjectPolicyCheckUC) Execute(context.Context) error { f.called++; return nil }

func TestGetProjectPolicyCheckCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustGetCfg(t)
	uc := &fakeProjectPolicyCheckUC{}
	ucs := defaultGetUCs()
	ucs.projectPolicyCheck = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policy-check", "-p", projectID.String()))
	require.Equal(t, 1, uc.called)
	require.Equal(t, projectID, cfg.ProjectId())
}
