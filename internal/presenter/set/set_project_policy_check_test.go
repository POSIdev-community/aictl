package set

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeSetPolicyCheckUC struct {
	called  int
	enabled bool
}

func (f *fakeSetPolicyCheckUC) Execute(_ context.Context, enabled bool) error {
	f.called++
	f.enabled = enabled
	return nil
}

func TestSetProjectPolicyCheckCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("true", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPolicyCheckUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policy-check", "true", "-p", projectID.String()))
		require.Equal(t, 1, uc.called)
		require.True(t, uc.enabled)
	})

	t.Run("false", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPolicyCheckUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policy-check", "false", "-p", projectID.String()))
		require.True(t, uc.called == 1 && !uc.enabled)
	})

	t.Run("invalid", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPolicyCheckUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, noopSetPoliciesUC{}, uc, &fakeSetExclusionsUC{})
		require.Error(t, cmdtest.Execute(t, root.Command, "project", "policy-check", "yes", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})
}
