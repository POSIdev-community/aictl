package policycheck

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

type fakeAI struct {
	body   string
	setRaw []byte
}

func (f *fakeAI) InitializeWithRetry(context.Context) error { return nil }

func (f *fakeAI) GetProjectPolicies(context.Context, uuid.UUID) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(f.body)), nil
}

func (f *fakeAI) SetProjectPolicies(_ context.Context, _ uuid.UUID, rawJSON []byte) error {
	f.setRaw = append([]byte(nil), rawJSON...)
	return nil
}

func TestUseCase_Execute_PreservesRules(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "", false, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.Nil)
	ai := &fakeAI{body: `{"checkSecurityPoliciesAccordance":false,"securityPolicies":"[{\"CountToActualize\":1}]"}`}

	uc, err := NewUseCase(ai, struct{}{}, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(context.Background(), true))

	var model struct {
		Check    bool   `json:"checkSecurityPoliciesAccordance"`
		Policies string `json:"securityPolicies"`
	}
	require.NoError(t, json.Unmarshal(ai.setRaw, &model))
	require.True(t, model.Check)
	require.JSONEq(t, `[{"CountToActualize":1}]`, model.Policies)
}
