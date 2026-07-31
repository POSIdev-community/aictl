package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeProjectSettingsUC struct {
	called int
	json   bool
}

func (f *fakeProjectSettingsUC) Execute(_ context.Context, jsonOutput bool) error {
	f.called++
	f.json = jsonOutput
	return nil
}

func TestGetProjectSettingsCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := mustGetCfg(t)
	uc := &fakeProjectSettingsUC{}
	ucs := defaultGetUCs()
	ucs.projectSettings = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "project", "settings", "-p", projectID.String(), "--json"))
	require.True(t, uc.json)
}
