package securitypolicies

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAndEncode(t *testing.T) {
	t.Parallel()

	t.Run("model_with_string", func(t *testing.T) {
		t.Parallel()
		m, err := Parse([]byte(`{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[{\"CountToActualize\":1}]"}`))
		require.NoError(t, err)
		require.True(t, m.CheckAccordance)
		require.JSONEq(t, `[{"CountToActualize":1}]`, m.Policies)

		out, err := Encode(m)
		require.NoError(t, err)
		require.JSONEq(t, `{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[{\"CountToActualize\":1}]"}`, string(out))
	})

	t.Run("bare_array", func(t *testing.T) {
		t.Parallel()
		m, err := Parse([]byte(`[{"CountToActualize":2}]`))
		require.NoError(t, err)
		require.False(t, m.CheckAccordance)
		require.JSONEq(t, `[{"CountToActualize":2}]`, m.Policies)
	})

	t.Run("null_policies", func(t *testing.T) {
		t.Parallel()
		m, err := Parse([]byte(`{"checkSecurityPoliciesAccordance":false,"securityPolicies":null}`))
		require.NoError(t, err)
		require.Equal(t, "[]", m.Policies)
	})
}
