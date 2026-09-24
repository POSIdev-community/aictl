package policycheck

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

type fakeCLI struct{ text string }

func (f *fakeCLI) ReturnText(_ context.Context, text string) { f.text = text }

func TestGetUseCase_Execute(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "", false, uuid.MustParse("11111111-1111-1111-1111-111111111111"), uuid.Nil)
	ai := &fakeAI{body: `{"checkSecurityPoliciesAccordance":true,"securityPolicies":"[]"}`}
	cli := &fakeCLI{}

	uc, err := NewUseCase(ai, cli, cfg)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(context.Background()))
	require.Equal(t, "true", cli.text)
}
