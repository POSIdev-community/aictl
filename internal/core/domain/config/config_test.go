package config_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func TestCACertPathMutualExclusion(t *testing.T) {
	t.Run("set_cacert_ok", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, cfg.SetCACertPath("/tmp/ca.pem"))
		require.Equal(t, "/tmp/ca.pem", cfg.CACertPath())
	})

	t.Run("set_cacert_rejects_empty", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.Error(t, cfg.SetCACertPath(""))
	})

	t.Run("set_cacert_rejects_when_tls_skip", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", true, uuid.Nil, uuid.Nil)
		require.Error(t, cfg.SetCACertPath("/tmp/ca.pem"))
	})

	t.Run("set_tls_skip_rejects_when_cacert", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, cfg.SetCACertPath("/tmp/ca.pem"))
		require.Error(t, cfg.SetTLSSkip(true))
		require.Equal(t, "/tmp/ca.pem", cfg.CACertPath())
		require.False(t, cfg.TLSSkip())
	})

	t.Run("validate_rejects_both", func(t *testing.T) {
		uri, err := config.NewUri("https://ai.example")
		require.NoError(t, err)
		cfg := config.NewConfig(uri, "token", true, uuid.Nil, uuid.Nil)
		require.Error(t, cfg.ApplyCACertPath("/tmp/ca.pem"))
	})

	t.Run("clear_cacert", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, cfg.SetCACertPath("/tmp/ca.pem"))
		cfg.ClearCACertPath()
		require.Empty(t, cfg.CACertPath())
	})
}
