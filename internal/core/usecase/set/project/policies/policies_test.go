package policies

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type fakeSetAI struct {
	getBody string
	setRaw  []byte
}

func (f *fakeSetAI) InitializeWithRetry(context.Context) error { return nil }

func (f *fakeSetAI) GetProjectPolicies(context.Context, uuid.UUID) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(f.getBody)), nil
}

func (f *fakeSetAI) SetProjectPolicies(_ context.Context, _ uuid.UUID, rawJSON []byte) error {
	f.setRaw = append([]byte(nil), rawJSON...)
	return nil
}

func TestUseCase_Execute_PreservesCheck(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "", false, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.Nil)
	ai := &fakeSetAI{getBody: `{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`}

	uc, err := NewUseCase(ai, struct{}{}, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(context.Background(), []byte(`[{"CountToActualize":1}]`), nil))

	var model struct {
		Check    bool   `json:"checkSecurityPoliciesAccordance"`
		Policies string `json:"securityPolicies"`
	}
	require.NoError(t, json.Unmarshal(ai.setRaw, &model))
	require.True(t, model.Check)
	require.JSONEq(t, `[{"CountToActualize":1}]`, model.Policies)
}

func TestUseCase_Execute_ExplicitCheck(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "", false, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.Nil)
	ai := &fakeSetAI{getBody: `{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`}
	check := false

	uc, err := NewUseCase(ai, struct{}{}, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(context.Background(), []byte(`[]`), &check))

	var model struct {
		Check bool `json:"checkSecurityPoliciesAccordance"`
	}
	require.NoError(t, json.Unmarshal(ai.setRaw, &model))
	require.False(t, model.Check)
}
