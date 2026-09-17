package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func TestFormatHelpers(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "512 B", formatIECSize(512))
	assert.Equal(t, "2.0 KiB", formatIECSize(2048))
	assert.Equal(t, "1.0 MiB", formatIECSize(1024*1024))
	assert.Equal(t, "abcdef01", shortHash("abcdef0123456789"))
	assert.Equal(t, "-", formatPackageUser(nil))
	assert.Equal(t, "ivan", formatPackageUser(&scafeeds.PackageUser{UserName: "ivan"}))
	assert.Equal(t, "ivan <a@b.c>", formatPackageUser(&scafeeds.PackageUser{UserName: "ivan", Email: "a@b.c"}))
	assert.Equal(t, "ivan <a@b.c> (tok)", formatPackageUser(&scafeeds.PackageUser{UserName: "ivan", Email: "a@b.c", TokenName: "tok"}))
	assert.Equal(t, "2026-09-11T18:42:00Z", formatRFC3339(time.Date(2026, 9, 11, 18, 42, 0, 0, time.UTC)))
}

func TestShowScaFeedsHeadersOnly(t *testing.T) {
	t.Parallel()

	var out, errBuf bytes.Buffer
	a := &Adapter{stdout: &out}
	a.ShowScaFeeds(testCtx(&out, &errBuf), nil)
	require.Contains(t, out.String(), "VERSION")
	require.Contains(t, out.String(), "STATUS")
	require.NotContains(t, out.String(), "TYPE")
}

func TestShowScaFeedsAlignedColumns(t *testing.T) {
	t.Parallel()

	var out, errBuf bytes.Buffer
	a := &Adapter{stdout: &out}
	pkgs := []scafeeds.Package{{
		PackageType: scafeeds.TypeScaFeeds,
		Version:     "1.0.159",
		Status:      scafeeds.StatusCurrent,
		FileName:    "AI.SCA.Feeds.1.0.159.zip",
		FileSize:    254_541_824,
		FileHash:    "423aa9bcdeadbeef",
		TriggeredBy: scafeeds.TriggerScheduled,
		UploadedAt:  time.Date(2026, 9, 2, 10, 58, 58, 0, time.UTC),
	}}
	a.ShowScaFeeds(testCtx(&out, &errBuf), pkgs)

	lines := bytes.Split(bytes.TrimRight(out.Bytes(), "\n"), []byte("\n"))
	require.Len(t, lines, 2)

	header := string(lines[0])
	row := string(lines[1])
	require.Contains(t, row, "-")
	require.NotContains(t, row, "n/a")
	require.NotContains(t, header, "TYPE")
	require.Equal(t, strings.Index(header, "VERSION"), strings.Index(row, "1.0.159"))
	require.Equal(t, strings.Index(header, "STATUS"), strings.Index(row, "current"))
	require.Equal(t, strings.Index(header, "HASH"), strings.Index(row, "423aa9bc"))
}
