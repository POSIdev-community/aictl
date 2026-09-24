package set

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

func resetSetProjectID() { projectIdFlag = "" }

type fakeSetExclusionsUC struct {
	called     int
	exclusions string
}

func (f *fakeSetExclusionsUC) Execute(_ context.Context, exclusions string) error {
	f.called++
	f.exclusions = exclusions
	return nil
}

type noopSetPoliciesUC struct{}

func (noopSetPoliciesUC) Execute(context.Context, []byte, *bool) error { return nil }

type noopSetPolicyCheckUC struct{}

func (noopSetPolicyCheckUC) Execute(context.Context, bool) error { return nil }

type noopSetSettingsUC struct{}

func (noopSetSettingsUC) Execute(context.Context, []byte) error { return nil }

func buildSetRoot(
	t *testing.T,
	settings UseCaseSetProjectSettings,
	policies UseCaseSetProjectPolicies,
	policyCheck UseCaseSetProjectPolicyCheck,
	exclusions UseCaseSetProjectExclusions,
) *CmdSet {
	t.Helper()
	cfg := cmdtest.MustCfg(t)
	preSet := NewPersistentPreRunESetCmd(cfg)
	preProj := NewPersistentPreRunESetProjectCmd(cfg, preSet)
	proj := NewSetProjectCmd(preProj,
		NewSetProjectSettingsCmd(settings),
		NewSetProjectPoliciesCmd(policies),
		NewSetProjectPolicyCheckCmd(policyCheck),
		NewSetProjectExclusionsCmd(exclusions),
	)
	return NewSetCmd(preSet, proj)
}

func TestSetProjectExclusionsCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("from_arg", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetExclusionsUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, noopSetPolicyCheckUC{}, uc)
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "exclusions", "*.tmp", "-p", projectID.String()))
		require.Equal(t, "*.tmp", uc.exclusions)
	})

	t.Run("from_file", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "exclusions.txt")
		require.NoError(t, os.WriteFile(path, []byte("node_modules/\n"), 0o644))
		uc := &fakeSetExclusionsUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, noopSetPolicyCheckUC{}, uc)
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "exclusions", "-f", path, "-p", projectID.String()))
		require.Equal(t, "node_modules/\n", uc.exclusions)
	})

	t.Run("stdin_arg_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetExclusionsUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, noopSetPolicyCheckUC{}, uc)
		cmdtest.WithStdin(t, "*.log\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "exclusions", "-", "-p", projectID.String()))
		})
		require.Equal(t, "*.log", uc.exclusions)
	})

	t.Run("stdin_file_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetExclusionsUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, noopSetPolicyCheckUC{}, uc)
		cmdtest.WithStdin(t, "vendor/\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "exclusions", "-f", "-", "-p", projectID.String()))
		})
		require.Equal(t, "vendor/\n", uc.exclusions)
	})
}
