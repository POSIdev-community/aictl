package common

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

func TestMatchesVersionRange(t *testing.T) {
	t.Parallel()

	ranges := []struct {
		name  string
		min   string
		max   string
		cases []struct {
			version string
			want    bool
		}
	}{
		{
			name: "v5_x",
			min:  "5.0.0",
			max:  "6.0.0",
			cases: []struct {
				version string
				want    bool
			}{
				{version: "4.9.9", want: false},
				{version: "5.0.0", want: true},
				{version: "5.4.0", want: true},
				{version: "5.99.0", want: true},
				{version: "6.0.0", want: false},
			},
		},
		{
			name: "v6_0",
			min:  "6.0.0",
			max:  "6.1.0",
			cases: []struct {
				version string
				want    bool
			}{
				{version: "5.99.0", want: false},
				{version: "6.0.0", want: true},
				{version: "6.0.5", want: true},
				{version: "6.1.0", want: false},
				{version: "7.0.0", want: false},
			},
		},
		{
			name: "v6_1",
			min:  "6.1.0",
			max:  "6.3.0",
			cases: []struct {
				version string
				want    bool
			}{
				{version: "6.0.0", want: false},
				{version: "6.0.5", want: false},
				{version: "6.1.0", want: true},
				{version: "6.2.0", want: true},
				{version: "6.2.9", want: true},
				{version: "6.3.0", want: false},
			},
		},
		{
			name: "v6_x",
			min:  "6.3.0",
			max:  "7.0.0",
			cases: []struct {
				version string
				want    bool
			}{
				{version: "6.2.9", want: false},
				{version: "6.3.0", want: true},
				{version: "6.9.0", want: true},
				{version: "7.0.0", want: false},
			},
		},
	}

	for _, r := range ranges {
		t.Run(r.name, func(t *testing.T) {
			t.Parallel()

			min, err := version.NewVersion(r.min)
			if err != nil {
				t.Fatalf("new min version: %v", err)
			}

			max, err := version.NewVersion(r.max)
			if err != nil {
				t.Fatalf("new max version: %v", err)
			}

			for _, tt := range r.cases {
				t.Run(tt.version, func(t *testing.T) {
					t.Parallel()

					ver, err := version.NewVersion(tt.version)
					if err != nil {
						t.Fatalf("new version: %v", err)
					}

					if got := MatchesVersionRange(ver, min, max); got != tt.want {
						t.Fatalf("MatchesVersionRange(%s, %s, %s) = %v, want %v", tt.version, r.min, r.max, got, tt.want)
					}
				})
			}
		})
	}
}
