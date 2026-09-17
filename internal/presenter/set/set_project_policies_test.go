package set

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeSetPoliciesUC struct {
	called  int
	rawJSON []byte
}

func (f *fakeSetPoliciesUC) Execute(_ context.Context, rawJSON []byte) error {
	f.called++
	f.rawJSON = append([]byte(nil), rawJSON...)
	return nil
}

func requirePoliciesPayload(t *testing.T, raw []byte, wantPolicies string, wantCheck bool) {
	t.Helper()
	var model struct {
		Check    bool   `json:"checkSecurityPoliciesAccordance"`
		Policies string `json:"securityPolicies"`
	}
	require.NoError(t, json.Unmarshal(raw, &model))
	require.Equal(t, wantCheck, model.Check)
	require.JSONEq(t, wantPolicies, model.Policies)
}

func TestSetProjectPoliciesCmd(t *testing.T) {
	t.Cleanup(resetSetProjectID)
	projectID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	t.Run("invalid_json", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.Error(t, cmdtest.Execute(t, root.Command, "project", "policies", "not-json", "-p", projectID.String()))
		require.Equal(t, 0, uc.called)
	})

	t.Run("from_arg_model", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies",
			`{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[]"}`, "-p", projectID.String()))
		requirePoliciesPayload(t, uc.rawJSON, `[]`, false)
	})

	t.Run("from_file_array", func(t *testing.T) {
		resetSetProjectID()
		path := filepath.Join(t.TempDir(), "policies.json")
		require.NoError(t, os.WriteFile(path, []byte(`[{"CountToActualize":1}]`), 0o644))
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", path, "-p", projectID.String()))
		requirePoliciesPayload(t, uc.rawJSON, `[{"CountToActualize":1}]`, false)
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
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", path, "-p", projectID.String()))

		var model struct {
			Check    bool   `json:"checkSecurityPoliciesAccordance"`
			Policies string `json:"securityPolicies"`
		}
		require.NoError(t, json.Unmarshal(uc.rawJSON, &model))
		require.False(t, model.Check)
		require.Contains(t, model.Policies, "// field name")
		require.Contains(t, model.Policies, `"Field": "VulnerabilityLevel"`)
	})

	t.Run("stdin_arg_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `[]`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-", "-p", projectID.String()))
		})
		requirePoliciesPayload(t, uc.rawJSON, `[]`, false)
	})

	t.Run("stdin_file_dash", func(t *testing.T) {
		resetSetProjectID()
		uc := &fakeSetPoliciesUC{}
		root := buildSetRoot(t, noopSetSettingsUC{}, uc, &fakeSetExclusionsUC{})
		cmdtest.WithStdin(t, `[{"CountToActualize":2}]`+"\n", func() {
			require.NoError(t, cmdtest.Execute(t, root.Command, "project", "policies", "-f", "-", "-p", projectID.String()))
		})
		requirePoliciesPayload(t, uc.rawJSON, `[{"CountToActualize":2}]`, false)
	})
}
