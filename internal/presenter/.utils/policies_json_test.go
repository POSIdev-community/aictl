package _utils

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePoliciesJSON(t *testing.T) {
	t.Parallel()

	t.Run("wraps_array_preserving_comments", func(t *testing.T) {
		t.Parallel()
		in := []byte(`[
    {
        "CountToActualize": 1,
        "Scopes": [
            {
                "Rules": [
                    {
                        "Field": "VulnerabilityLevel", // field name
                        "Value": "High",
                        "IsRegex": false
                    }
                ]
            }
        ]
    }
]`)
		out, err := NormalizePoliciesJSON(in)
		require.NoError(t, err)

		var model struct {
			Check    bool   `json:"checkSecurityPoliciesAccordance"`
			Policies string `json:"securityPolicies"`
		}
		require.NoError(t, json.Unmarshal(out, &model))
		require.False(t, model.Check)
		require.Contains(t, model.Policies, "// field name")
		require.Contains(t, model.Policies, `"Field": "VulnerabilityLevel"`)
	})

	t.Run("passes_through_model_with_string", func(t *testing.T) {
		t.Parallel()
		in := []byte(`{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`)
		out, err := NormalizePoliciesJSON(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`, string(out))
	})

	t.Run("stringifies_array_securityPolicies", func(t *testing.T) {
		t.Parallel()
		in := []byte(`{"checkSecurityPoliciesAccordance":true,"securityPolicies":[{"CountToActualize":1}]}`)
		out, err := NormalizePoliciesJSON(in)
		require.NoError(t, err)

		var model struct {
			Check    bool   `json:"checkSecurityPoliciesAccordance"`
			Policies string `json:"securityPolicies"`
		}
		require.NoError(t, json.Unmarshal(out, &model))
		require.True(t, model.Check)
		require.JSONEq(t, `[{"CountToActualize":1}]`, model.Policies)
	})

	t.Run("wraps_bare_object", func(t *testing.T) {
		t.Parallel()
		out, err := NormalizePoliciesJSON([]byte(`{"CountToActualize":1}`))
		require.NoError(t, err)

		var model struct {
			Policies string `json:"securityPolicies"`
		}
		require.NoError(t, json.Unmarshal(out, &model))
		require.JSONEq(t, `[{"CountToActualize":1}]`, model.Policies)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		_, err := NormalizePoliciesJSON([]byte(`not-json`))
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "array or object"))
	})
}
