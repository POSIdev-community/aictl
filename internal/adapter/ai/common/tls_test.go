package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func writeTempPEM(t *testing.T) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "aictl-test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "ca.pem")
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: der}))
	require.NoError(t, f.Close())

	return path
}

func TestNewTLSConfig(t *testing.T) {
	t.Run("default_nil", func(t *testing.T) {
		cfg, err := NewTLSConfig(false, "", nil)
		require.NoError(t, err)
		require.Nil(t, cfg)
	})

	t.Run("tls_skip", func(t *testing.T) {
		cfg, err := NewTLSConfig(true, "", nil)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.True(t, cfg.InsecureSkipVerify)
	})

	t.Run("tls_skip_preserves_next_protos", func(t *testing.T) {
		base := &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
		cfg, err := NewTLSConfig(true, "", base)
		require.NoError(t, err)
		require.Equal(t, []string{"h2", "http/1.1"}, cfg.NextProtos)
	})

	t.Run("mutual_exclusion", func(t *testing.T) {
		_, err := NewTLSConfig(true, "/tmp/ca.pem", nil)
		require.Error(t, err)
	})

	t.Run("append_pem", func(t *testing.T) {
		path := writeTempPEM(t)
		cfg, err := NewTLSConfig(false, path, nil)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.RootCAs)
		require.False(t, cfg.InsecureSkipVerify)
	})

	t.Run("missing_file", func(t *testing.T) {
		_, err := NewTLSConfig(false, filepath.Join(t.TempDir(), "missing.pem"), nil)
		require.Error(t, err)
	})

	t.Run("empty_pem", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.pem")
		require.NoError(t, os.WriteFile(path, []byte("not a cert\n"), 0o600))
		_, err := NewTLSConfig(false, path, nil)
		require.Error(t, err)
	})
}

func TestNewHTTPTransport(t *testing.T) {
	tr, err := NewHTTPTransport(false, "")
	require.NoError(t, err)
	require.NotNil(t, tr)
	// Must keep HTTP/2 from DefaultTransport (zero Transport.Clone() breaks h2 ALPN).
	require.True(t, tr.ForceAttemptHTTP2)

	path := writeTempPEM(t)
	tr, err = NewHTTPTransport(false, path)
	require.NoError(t, err)
	require.NotNil(t, tr.TLSClientConfig)
	require.NotNil(t, tr.TLSClientConfig.RootCAs)
	require.True(t, tr.ForceAttemptHTTP2)
	require.Contains(t, tr.TLSClientConfig.NextProtos, "h2")
}

func TestNewHTTPTransport_CloneSafeForHTTPS(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	caPath := filepath.Join(t.TempDir(), "server-ca.pem")
	require.NoError(t, os.WriteFile(caPath, certPEM, 0o600))

	tr, err := NewHTTPTransport(false, caPath)
	require.NoError(t, err)
	require.True(t, tr.ForceAttemptHTTP2)

	cloned := tr.Clone()
	require.True(t, cloned.ForceAttemptHTTP2, "Clone must preserve HTTP/2 settings")

	client := &http.Client{Transport: cloned}
	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
