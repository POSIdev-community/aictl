package show

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type fakeCFG struct {
	human, json, yaml string
	called            string
}

func (f *fakeCFG) String(_ *config.Config) (string, error) {
	f.called = "string"
	return f.human, nil
}

func (f *fakeCFG) StringJson(_ *config.Config) (string, error) {
	f.called = "json"
	return f.json, nil
}

func (f *fakeCFG) StringYaml(_ *config.Config) (string, error) {
	f.called = "yaml"
	return f.yaml, nil
}

type stubCLI struct {
	out string
}

func (s *stubCLI) ReturnText(_ context.Context, text string) {
	s.out = text
}

func TestUseCase_Execute_routesFormats(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig(config.Uri{}, "token", false, uuid.Nil, uuid.Nil)

	tests := []struct {
		name       string
		json, yaml bool
		wantCalled string
		wantOut    string
	}{
		{name: "default", wantCalled: "string", wantOut: "human-out"},
		{name: "json", json: true, wantCalled: "json", wantOut: "json-out"},
		{name: "yaml", yaml: true, wantCalled: "yaml", wantOut: "yaml-out"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfgAdapter := &fakeCFG{human: "human-out", json: "json-out", yaml: "yaml-out"}
			cli := &stubCLI{}
			uc, err := NewUseCase(cfgAdapter, cli, cfg)
			require.NoError(t, err)

			require.NoError(t, uc.Execute(t.Context(), tt.json, tt.yaml))
			require.Equal(t, tt.wantCalled, cfgAdapter.called)
			require.Equal(t, tt.wantOut, cli.out)
		})
	}
}

func TestUseCase_Execute_jsonAndYamlRejected(t *testing.T) {
	t.Parallel()

	cfgAdapter := &fakeCFG{}
	cli := &stubCLI{}
	uc, err := NewUseCase(cfgAdapter, cli, config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil))
	require.NoError(t, err)

	err = uc.Execute(t.Context(), true, true)
	require.Error(t, err)
	require.Empty(t, cfgAdapter.called)
	require.Empty(t, cli.out)
}
