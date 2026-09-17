package config

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	domainconfig "github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func fullyPopulatedConfig(t *testing.T) *domainconfig.Config {
	t.Helper()

	uri, err := domainconfig.NewUri("https://example.com/ai")
	require.NoError(t, err)

	projectID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174011")
	branchID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	cfg := domainconfig.NewConfig(uri, "super-secret-token", false, projectID, branchID)
	require.NoError(t, cfg.SetCACertPath("/etc/ssl/certs/custom-ca.pem"))

	return cfg
}

func TestAdapter_String_masksOnlyToken(t *testing.T) {
	t.Parallel()

	cfg := fullyPopulatedConfig(t)
	adapter := NewContextAdapter()

	out, err := adapter.String(cfg)
	require.NoError(t, err)

	require.Contains(t, out, "token: "+maskedToken)
	require.NotContains(t, out, "super-secret-token")

	require.Contains(t, out, "uri: https://example.com/ai")
	require.Contains(t, out, "projectId: 123e4567-e89b-12d3-a456-426614174011")
	require.Contains(t, out, "branchId: 123e4567-e89b-12d3-a456-426614174000")
	require.Contains(t, out, "ca-cert: /etc/ssl/certs/custom-ca.pem")
	require.Contains(t, out, "tls-skip: false")
	require.NotContains(t, out, unsetValue)
}

func TestAdapter_String_unsetFieldsStayUnsetMarker(t *testing.T) {
	t.Parallel()

	cfg := domainconfig.NewConfig(domainconfig.Uri{}, "", false, uuid.Nil, uuid.Nil)
	adapter := NewContextAdapter()

	out, err := adapter.String(cfg)
	require.NoError(t, err)

	require.Contains(t, out, "uri: "+unsetValue)
	require.Contains(t, out, "token: "+unsetValue)
	require.Contains(t, out, "projectId: "+unsetValue)
	require.Contains(t, out, "branchId: "+unsetValue)
	require.Contains(t, out, "ca-cert: "+unsetValue)
	require.NotContains(t, out, maskedToken)
}

func TestAdapter_StringJson_rawTokenAndNullUnset(t *testing.T) {
	t.Parallel()

	cfg := fullyPopulatedConfig(t)
	adapter := NewContextAdapter()

	out, err := adapter.StringJson(cfg)
	require.NoError(t, err)
	require.NotContains(t, out, maskedToken)
	require.NotContains(t, out, unsetValue)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &decoded))
	require.Equal(t, "super-secret-token", decoded["token"])
	require.Equal(t, "https://example.com/ai", decoded["uri"])
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174011", decoded["projectId"])
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", decoded["branchId"])
	require.Equal(t, "/etc/ssl/certs/custom-ca.pem", decoded["caCert"])
	require.Equal(t, false, decoded["tlsSkip"])
}

func TestAdapter_StringJson_nullForUnsetFields(t *testing.T) {
	t.Parallel()

	uri, err := domainconfig.NewUri("https://example.com")
	require.NoError(t, err)

	cfg := domainconfig.NewConfig(uri, "token-value", true, uuid.Nil, uuid.Nil)
	adapter := NewContextAdapter()

	out, err := adapter.StringJson(cfg)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &decoded))
	require.Equal(t, "token-value", decoded["token"])
	require.Equal(t, "https://example.com", decoded["uri"])
	require.Equal(t, true, decoded["tlsSkip"])
	require.Nil(t, decoded["projectId"])
	require.Nil(t, decoded["branchId"])
	require.Nil(t, decoded["caCert"])
}

func TestAdapter_StringYaml_rawTokenAndNullUnset(t *testing.T) {
	t.Parallel()

	cfg := fullyPopulatedConfig(t)
	adapter := NewContextAdapter()

	out, err := adapter.StringYaml(cfg)
	require.NoError(t, err)
	require.NotContains(t, out, maskedToken)
	require.NotContains(t, out, unsetValue)

	var decoded fileConfigExport
	require.NoError(t, yaml.Unmarshal([]byte(out), &decoded))
	require.NotNil(t, decoded.Token)
	require.Equal(t, "super-secret-token", *decoded.Token)
	require.NotNil(t, decoded.Uri)
	require.Equal(t, "https://example.com/ai", *decoded.Uri)
	require.NotNil(t, decoded.ProjectId)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174011", *decoded.ProjectId)
	require.NotNil(t, decoded.BranchId)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", *decoded.BranchId)
	require.NotNil(t, decoded.CACert)
	require.Equal(t, "/etc/ssl/certs/custom-ca.pem", *decoded.CACert)
	require.False(t, decoded.TLSSkip)
}

func TestAdapter_StringYaml_nullForUnsetFields(t *testing.T) {
	t.Parallel()

	cfg := domainconfig.NewConfig(domainconfig.Uri{}, "", false, uuid.Nil, uuid.Nil)
	adapter := NewContextAdapter()

	out, err := adapter.StringYaml(cfg)
	require.NoError(t, err)
	require.NotContains(t, out, unsetValue)

	var decoded fileConfigExport
	require.NoError(t, yaml.Unmarshal([]byte(out), &decoded))
	require.Nil(t, decoded.Uri)
	require.Nil(t, decoded.Token)
	require.Nil(t, decoded.ProjectId)
	require.Nil(t, decoded.BranchId)
	require.Nil(t, decoded.CACert)
	require.False(t, decoded.TLSSkip)
}

func TestStoreContext_keepsRawToken(t *testing.T) {
	t.Parallel()

	uri, err := domainconfig.NewUri("https://example.com")
	require.NoError(t, err)

	cfg := domainconfig.NewConfig(uri, "super-secret-token", false, uuid.Nil, uuid.Nil)
	fileCfg := fileConfigFromDomainConfig(cfg)

	require.Equal(t, "super-secret-token", fileCfg.Token)
	require.Empty(t, fileCfg.ProjectId)
}
