package scafeeds

import "strings"

// SanitizeFileName keeps [A-Za-z0-9._-], replaces everything else with '_'.
func SanitizeFileName(name string) string {
	if name == "" {
		return name
	}

	var b strings.Builder
	b.Grow(len(name))

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	return b.String()
}

// DefaultArchiveName returns sca-feeds-<sanitized-version>.zip.
func DefaultArchiveName(version string) string {
	return "sca-feeds-" + SanitizeFileName(version) + ".zip"
}
