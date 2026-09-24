package _utils

import (
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
		require.Nil(t, out.Check)
		require.Contains(t, string(out.PoliciesJSON), "// field name")
		require.Contains(t, string(out.PoliciesJSON), `"Field": "VulnerabilityLevel"`)
	})

	t.Run("passes_through_model_with_string", func(t *testing.T) {
		t.Parallel()
		in := []byte(`{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`)
		out, err := NormalizePoliciesJSON(in)
		require.NoError(t, err)
		require.NotNil(t, out.Check)
		require.True(t, *out.Check)
		require.Equal(t, "[]", string(out.PoliciesJSON))
	})

	t.Run("stringifies_array_securityPolicies", func(t *testing.T) {
		t.Parallel()
		in := []byte(`{"checkSecurityPoliciesAccordance":true,"securityPolicies":[{"CountToActualize":1}]}`)
		out, err := NormalizePoliciesJSON(in)
		require.NoError(t, err)
		require.NotNil(t, out.Check)
		require.True(t, *out.Check)
		require.JSONEq(t, `[{"CountToActualize":1}]`, string(out.PoliciesJSON))
	})

	t.Run("wraps_bare_object", func(t *testing.T) {
		t.Parallel()
		out, err := NormalizePoliciesJSON([]byte(`{"CountToActualize":1}`))
		require.NoError(t, err)
		require.Nil(t, out.Check)
		require.JSONEq(t, `[{"CountToActualize":1}]`, string(out.PoliciesJSON))
	})

	t.Run("model_without_check_preserves_nil", func(t *testing.T) {
		t.Parallel()
		out, err := NormalizePoliciesJSON([]byte(`{"securityPolicies":[]}`))
		require.NoError(t, err)
		require.Nil(t, out.Check)
		require.Equal(t, "[]", string(out.PoliciesJSON))
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		_, err := NormalizePoliciesJSON([]byte(`not-json`))
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "array or object"))
	})
}
