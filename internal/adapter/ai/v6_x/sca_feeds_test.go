package v6_x

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	clientai "github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

func newTestClient6x(t *testing.T, h http.Handler) *ClientAI6x {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	base := common.NewBaseClient()
	base.HttpClient = srv.Client()
	c := NewAiClient(base)
	api, err := clientai.NewClientWithResponses(srv.URL, clientai.WithHTTPClient(srv.Client()))
	require.NoError(t, err)
	c.ClientWithResponses = api

	return c
}

func sampleAPIPackage(id int32, version string, status clientai.PackageStatus, fileName string) clientai.Package {
	now := time.Date(2026, 9, 11, 18, 42, 0, 0, time.UTC)
	uid := openapi_types.UUID(uuid.MustParse("11111111-1111-1111-1111-111111111111"))
	user := clientai.PackageUser{Id: uid, UserName: "ivan", Email: "ivan@example.com", TokenName: "ci"}

	return clientai.Package{
		Id:             id,
		PackageType:    clientai.ScaFeeds,
		Version:        version,
		Status:         status,
		FileName:       fileName,
		FileSize:       2048,
		FileHash:       "abcdef0123456789",
		TriggeredBy:    clientai.PackageTriggerManual,
		UploadedAt:     now,
		UploadedBy:     &user,
		LastModifiedAt: now,
		LastModifiedBy: &user,
	}
}

func TestGetScaFeeds(t *testing.T) {
	t.Parallel()

	pkgs := []clientai.Package{
		sampleAPIPackage(1, "1.0", clientai.Current, "a.zip"),
		sampleAPIPackage(2, "0.9", clientai.Active, "b.zip"),
	}

	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/packages", r.URL.Path)
		assert.Equal(t, "sca_feeds", r.URL.Query().Get("type"))
		assert.Empty(t, r.URL.Query().Get("status"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pkgs)
	}))

	got, err := client.GetScaFeeds(t.Context(), []scafeeds.Status{scafeeds.StatusCurrent})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "1.0", got[0].Version)
	assert.Equal(t, scafeeds.StatusCurrent, got[0].Status)
}

func TestDownloadScaFeeds(t *testing.T) {
	t.Parallel()

	list := []clientai.Package{sampleAPIPackage(7, "1.2.3", clientai.Active, "feeds.zip")}
	payload := bytes.Repeat([]byte("Z"), 2048) // matches sample FileSize
	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages/7/content":
			w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))

	var errBuf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		}),
		zapcore.AddSync(&errBuf),
		zapcore.ErrorLevel,
	)
	ctx := logger.ContextWithLogger(t.Context(), logger.Wrap(zap.New(core)))

	rc, name, err := client.DownloadScaFeeds(ctx, "1.2.3")
	require.NoError(t, err)
	defer func() { _ = rc.Close() }()
	assert.Equal(t, "feeds.zip", name)
	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, payload, data)

	out := errBuf.String()
	require.Contains(t, out, "downloading sca feeds: 0%")
	require.Contains(t, out, "downloading sca feeds: 100%")
}

func TestDownloadScaFeeds_QuietNoProgress(t *testing.T) {
	t.Parallel()

	list := []clientai.Package{sampleAPIPackage(7, "1.2.3", clientai.Active, "feeds.zip")}
	payload := bytes.Repeat([]byte("Z"), 2048)
	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages/7/content":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))

	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		}),
		zapcore.AddSync(&buf),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		}),
	)
	ctx := logger.ContextWithLogger(t.Context(), logger.Wrap(zap.New(core)))

	rc, _, err := client.DownloadScaFeeds(ctx, "1.2.3")
	require.NoError(t, err)
	defer func() { _ = rc.Close() }()
	_, err = io.Copy(io.Discard, rc)
	require.NoError(t, err)
	require.NotContains(t, buf.String(), "downloading sca feeds")
}

func TestDownloadScaFeeds_DefaultFileNameAndZeroSizeProgress(t *testing.T) {
	t.Parallel()

	pkg := sampleAPIPackage(8, "2.0", clientai.Current, "")
	pkg.FileSize = 0
	payload := bytes.Repeat([]byte("A"), 1024)

	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]clientai.Package{pkg})
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages/8/content":
			w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(payload)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))

	var errBuf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		}),
		zapcore.AddSync(&errBuf),
		zapcore.ErrorLevel,
	)
	ctx := logger.ContextWithLogger(t.Context(), logger.Wrap(zap.New(core)))

	rc, name, err := client.DownloadScaFeeds(ctx, "2.0")
	require.NoError(t, err)
	defer func() { _ = rc.Close() }()
	assert.Equal(t, "sca-feeds-2.0.zip", name)

	data, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, payload, data)

	out := errBuf.String()
	require.Contains(t, out, "downloading sca feeds: 0%")
	require.Contains(t, out, "downloading sca feeds: 100%")
}

func TestDownloadScaFeeds_HTTPError(t *testing.T) {
	t.Parallel()

	list := []clientai.Package{sampleAPIPackage(7, "1.2.3", clientai.Active, "feeds.zip")}
	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages/7/content":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("missing"))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))

	_, _, err := client.DownloadScaFeeds(t.Context(), "1.2.3")
	require.Error(t, err)
}

func TestDownloadScaFeedsNotFound(t *testing.T) {
	t.Parallel()

	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]clientai.Package{})
	}))

	_, _, err := client.DownloadScaFeeds(t.Context(), "missing")
	require.Error(t, err)
	var ve *validation.Error
	require.ErrorAs(t, err, &ve)
}

func TestRollbackScaFeeds(t *testing.T) {
	t.Parallel()

	current := sampleAPIPackage(10, "2.0", clientai.Current, "c.zip")
	rolled := sampleAPIPackage(9, "1.0", clientai.Current, "p.zip")

	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/packages":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]clientai.Package{current})
		case r.Method == http.MethodPost && r.URL.Path == "/api/packages/10/rollback":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(rolled)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))

	pkg, err := client.RollbackScaFeeds(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "1.0", pkg.Version)
}

func TestRollbackScaFeedsNoCurrent(t *testing.T) {
	t.Parallel()

	client := newTestClient6x(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]clientai.Package{
			sampleAPIPackage(1, "1.0", clientai.Active, "a.zip"),
		})
	}))

	_, err := client.RollbackScaFeeds(t.Context())
	require.Error(t, err)
	assert.Contains(t, err.Error(), scafeeds.ErrNoCurrent)
}
