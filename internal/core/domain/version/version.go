package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyVersion     = errors.New("empty version string")
	ErrInvalidComponent = errors.New("invalid version component")
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Build uint64

	components int
}

// NewVersion parses a version string like "1", "1.2", "1.2.3", or "1.2.3.4".
// Returns error if any component is negative or not a base-10 integer.
func NewVersion(s string) (Version, error) {
	if s == "" {
		return Version{}, ErrEmptyVersion
	}

	parts := strings.Split(s, ".")
	if len(parts) == 0 || len(parts) > 4 {
		return Version{}, fmt.Errorf("expected 1–4 dot-separated components, got %d", len(parts))
	}

	v := Version{components: len(parts)}
	for i, part := range parts {
		// Trim surrounding whitespace (robustness)
		part = strings.TrimSpace(part)
		if part == "" {
			return Version{}, fmt.Errorf("empty component at position %d", i+1)
		}

		num, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return Version{}, fmt.Errorf("%w at position %d: %q", ErrInvalidComponent, i+1, part)
		}

		switch i {
		case 0:
			v.Major = num
		case 1:
			v.Minor = num
		case 2:
			v.Patch = num
		case 3:
			v.Build = num
		}
	}

	return v, nil
}

// String returns canonical form, omitting trailing zero components.
// Examples:
//
//	"1" → "1"
//	"1.0" → "1.0"
//	"1.2.0" → "1.2.0"
//	"1.2.3.0" → "1.2.3.0"
//
// Unlike some implementations, we preserve explicit zeros to avoid ambiguity.
// (If you prefer minimal form like "1", change logic accordingly.)
func (v *Version) String() string {
	switch v.components {
	case 1:
		return strconv.FormatUint(v.Major, 10)
	case 2:
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	case 3:
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	default:
		return fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Build)
	}
}

func (v *Version) CompareVersion(other *Version) bool {
	var components int
	if v.components > other.components {
		components = other.components
	} else {
		components = v.components
	}

	switch components {
	case 1:
		return v.Major == other.Major
	case 2:
		return v.Major == other.Major && v.Minor == other.Minor
	case 3:
		return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
	default:
		return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch && v.Build == other.Build
	}
}

// Compare returns:
//
//	-1 if v < other
//	 0 if v == other
//	+1 if v > other
func (v *Version) Compare(other *Version) int {
	if diff := int64(v.Major) - int64(other.Major); diff != 0 {
		return sign(diff)
	}
	if diff := int64(v.Minor) - int64(other.Minor); diff != 0 {
		return sign(diff)
	}
	if diff := int64(v.Patch) - int64(other.Patch); diff != 0 {
		return sign(diff)
	}
	if diff := int64(v.Build) - int64(other.Build); diff != 0 {
		return sign(diff)
	}
	return 0
}

func sign(x int64) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// Less reports whether v < other.
func (v *Version) Less(other *Version) bool { return v.Compare(other) < 0 }

// Greater reports whether v > other.
func (v *Version) Greater(other *Version) bool { return v.Compare(other) > 0 }
