package update

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeUpdateLanguagesUC struct {
	called int
	err    error
	onCall func()
}

func (f *fakeUpdateLanguagesUC) Execute(context.Context) error {
	f.called++
	if f.onCall != nil {
		f.onCall()
	}
	return f.err
}

func TestUpdateProjectLanguagesCmd(t *testing.T) {
	t.Cleanup(resetUpdateFlags)
	resetUpdateFlags()
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cfg := cmdtest.MustCfg(t)
	uc := &fakeUpdateLanguagesUC{}
	projectCmd := NewUpdateProjectCmd(
		NewPersistentPreRunEUpdateProjectCmd(cfg, NewPersistentPreRunEUpdateCmd(cfg)),
		NewUpdateProjectSettingsCmd(noopUpdateSettingsUC{}),
		NewUpdateProjectLanguagesCmd(uc),
	)
	root := NewUpdateCmd(cfg, NewUpdateSourcesCmd(cfg, noopUpdateSourcesUC{}, noopUpdateLanguagesUC{}), NewUpdateSbomCmd(cfg, noopUpdateSbomUC{}), NewUpdateScaFeedsCmd(noopUpdateScaFeedsUC{}, noopRollbackScaFeedsUC{}), projectCmd)
	require.NoError(t, cmdtest.Execute(t, root.Command, "project", "languages", "-p", projectID.String()))
	require.Equal(t, 1, uc.called)
	require.Equal(t, projectID, cfg.ProjectId())
}
