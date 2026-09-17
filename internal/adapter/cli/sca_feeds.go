package cli

import (
	"fmt"
	"time"

	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
)

func formatIECSize(size int64) string {
	const (
		unit = 1024
	)

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	value := float64(size) / float64(div)
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", value, units[exp])
}

func shortHash(hash string) string {
	if len(hash) <= 8 {
		return hash
	}

	return hash[:8]
}

func formatRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339)
}

func formatPackageUser(u *scafeeds.PackageUser) string {
	if u == nil {
		return "-"
	}

	if u.Email == "" {
		name := u.UserName
		if u.TokenName != "" {
			return fmt.Sprintf("%s (%s)", name, u.TokenName)
		}

		return name
	}

	cell := fmt.Sprintf("%s <%s>", u.UserName, u.Email)
	if u.TokenName != "" {
		return fmt.Sprintf("%s (%s)", cell, u.TokenName)
	}

	return cell
}
