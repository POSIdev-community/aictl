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

	t.Run("validate_rejects_both", func(t *testing.T) {
		uri, err := config.NewUri("https://ai.example")
		require.NoError(t, err)
		cfg := config.NewConfig(uri, "token", false, uuid.Nil, uuid.Nil)
		require.NoError(t, cfg.SetCACertPath("/tmp/ca.pem"))
		cfg.SetTLSSkip(true)
		require.Error(t, cfg.Validate())
	})

	t.Run("clear_cacert", func(t *testing.T) {
		cfg := config.NewConfig(config.Uri{}, "", false, uuid.Nil, uuid.Nil)
		require.NoError(t, cfg.SetCACertPath("/tmp/ca.pem"))
		cfg.ClearCACertPath()
		require.Empty(t, cfg.CACertPath())
	})
}
