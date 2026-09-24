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

type fakeSetPoliciesUC struct {
	called       int
	policiesJSON []byte
	check        *bool
}

func (f *fakeSetPoliciesUC) Execute(_ context.Context, policiesJSON []byte, check *bool) error {
	f.called++
	f.policiesJSON = append([]byte(nil), policiesJSON...)
	f.check = check
	return nil
}

func TestSetProjectPoliciesCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("invalid_json", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		require.Error(t, cmdtest.Execute(t, root.Command, "project", "policies", "not-json", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})

	t.Run("from_arg_model", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies",
			`{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[]"}`, "-p", projectID.String()))
		require.Equal(t, "[]", string(uc.policiesJSON))
		require.NotNil(t, uc.check)
		require.False(t, *uc.check)
	})

	t.Run("from_file_array", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "policies.json")
		require.NoError(t, os.WriteFile(path, []byte(`[{"CountToActualize":1}]`), 0o644))
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", path, "-p", projectID.String()))
		require.JSONEq(t, `[{"CountToActualize":1}]`, string(uc.policiesJSON))
		require.Nil(t, uc.check)
	})

	t.Run("from_file_with_comments", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "policies.json")
		content := `[
    {
        "CountToActualize": 1,
        "Scopes": [
            {
                "Rules": [
                    {
                        "Field": "VulnerabilityLevel", // field name
                        "Value": "High",  // field value, case insensitive
                        "IsRegex": false  // whether to use regular expressions
                    }
                ]
            }
        ]
    }
]`
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", path, "-p", projectID.String()))
		require.Nil(t, uc.check)
		require.Contains(t, string(uc.policiesJSON), "// field name")
		require.Contains(t, string(uc.policiesJSON), `"Field": "VulnerabilityLevel"`)
	})

	t.Run("stdin_arg_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `[]`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-", "-p", projectID.String()))
		})
		require.Equal(t, "[]", string(uc.policiesJSON))
		require.Nil(t, uc.check)
	})

	t.Run("stdin_file_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, noopSetPolicyCheckUC{}, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `[{"CountToActualize":2}]`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", "-", "-p", projectID.String()))
		})
		require.JSONEq(t, `[{"CountToActualize":2}]`, string(uc.policiesJSON))
		require.Nil(t, uc.check)
	})
}
