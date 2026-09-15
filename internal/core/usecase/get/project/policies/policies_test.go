package policies

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type fakeAI struct {
	body string
	err  error
}

func (f *fakeAI) InitializeWithRetry(context.Context) error { return nil }

func (f *fakeAI) GetProjectPolicies(context.Context, uuid.UUID) (io.ReadCloser, error) {
	if f.err != nil {
		return nil, f.err
	}

	return io.NopCloser(strings.NewReader(f.body)), nil
}

type fakeCLI struct {
	text string
}

func (f *fakeCLI) ReturnText(_ context.Context, text string) {
	f.text = text
}

func TestFormatProjectPolicies(t *testing.T) {
	t.Parallel()

	t.Run("model_with_string", func(t *testing.T) {
		t.Parallel()
		in := `{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[{\"CountToActualize\":1}]"}`
		out, err := formatProjectPolicies([]byte(in))
		require.NoError(t, err)
		require.Equal(t, "[\n    {\n        \"CountToActualize\": 1\n    }\n]\n", out)
	})

	t.Run("model_with_array", func(t *testing.T) {
		t.Parallel()
		in := `{"checkSecurityPoliciesAccordance":true,"securityPolicies":[{"CountToActualize":2}]}`
		out, err := formatProjectPolicies([]byte(in))
		require.NoError(t, err)
		require.Equal(t, "[\n    {\n        \"CountToActualize\": 2\n    }\n]\n", out)
	})

	t.Run("null_policies", func(t *testing.T) {
		t.Parallel()
		out, err := formatProjectPolicies([]byte(`{"checkSecurityPoliciesAccordance":false,"securityPolicies":null}`))
		require.NoError(t, err)
		require.Equal(t, "[]\n", out)
	})

	t.Run("preserves_comments_in_string", func(t *testing.T) {
		t.Parallel()
		in := `{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[\n  {\"Field\": \"High\" // note\n}\n]"}`
		out, err := formatProjectPolicies([]byte(in))
		require.NoError(t, err)
		require.Contains(t, out, "// note")
		require.Contains(t, out, `"Field": "High"`)
	})
}

func TestUseCase_Execute_FormatsPolicies(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "", false, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.Nil)
	ai := &fakeAI{body: `{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[{\"CountToActualize\":1,\"Scopes\":[]}]"}`}
	cli := &fakeCLI{}

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(context.Background()))
	require.Equal(t, "[\n    {\n        \"CountToActualize\": 1,\n        \"Scopes\": []\n    }\n]\n", cli.text)
}
