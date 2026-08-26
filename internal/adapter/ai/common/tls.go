package common

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// NewHTTPTransport builds an HTTP transport with optional TLS skip or custom CA PEM.
// When caCertPath is set, certificates from the PEM file are appended to the system trust pool.
// tlsSkip and caCertPath are mutually exclusive.
//
// The transport is cloned from http.DefaultTransport so HTTP/2 stays enabled.
// A zero &http.Transport{}.Clone() disables HTTP/2 and breaks servers that negotiate h2.
func NewHTTPTransport(tlsSkip bool, caCertPath string) (*http.Transport, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	tlsCfg, err := NewTLSConfig(tlsSkip, caCertPath, transport.TLSClientConfig)
	if err != nil {
		return nil, err
	}
	if tlsCfg != nil {
		transport.TLSClientConfig = tlsCfg
	}

	return transport, nil
}

// NewTLSConfig returns a tls.Config for the given options.
// nil means the caller should keep the transport's existing TLS settings.
// base is optional (e.g. from DefaultTransport) and is used to preserve ALPN/NextProtos.
func NewTLSConfig(tlsSkip bool, caCertPath string, base *tls.Config) (*tls.Config, error) {
	if tlsSkip && caCertPath != "" {
		return nil, fmt.Errorf("cannot use cacert together with tls-skip")
	}

	if !tlsSkip && caCertPath == "" {
		return nil, nil
	}

	var cfg *tls.Config
	if base != nil {
		cfg = base.Clone()
	} else {
		cfg = &tls.Config{}
	}

	if tlsSkip {
		cfg.InsecureSkipVerify = true //nolint:gosec // intentional when user passes --tls-skip

		return cfg, nil
	}

	pemData, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("read cacert %q: %w", caCertPath, err)
	}

	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}

	if ok := pool.AppendCertsFromPEM(pemData); !ok {
		return nil, fmt.Errorf("cacert %q: no PEM certificates found", caCertPath)
	}

	cfg.RootCAs = pool
	cfg.InsecureSkipVerify = false

	return cfg, nil
}
