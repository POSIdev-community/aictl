package ai

import (
	"testing"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v5_x"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_0"
	"github.com/POSIdev-community/aictl/internal/adapter/ai/v6_x"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

var (
	_ ClientAi = (*v5_x.ClientAI5x)(nil)
	_ ClientAi = (*v6_0.ClientAI60)(nil)
	_ ClientAi = (*v6_x.ClientAI61)(nil)
)

func TestVersionRangeInitializerBounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		version  string
		min      string
		max      string
		expected bool
	}{
		{name: "5.0 matches v5_x", version: "5.0.0", min: "5.0.0", max: "6.0.0", expected: true},
		{name: "5.4 matches v5_x", version: "5.4.0", min: "5.0.0", max: "6.0.0", expected: true},
		{name: "4.9 below v5_x", version: "4.9.9", min: "5.0.0", max: "6.0.0", expected: false},
		{name: "6.0 matches v6_0", version: "6.0.0", min: "6.0.0", max: "6.1.0", expected: true},
		{name: "6.1 matches v6_x", version: "6.1.0", min: "6.1.0", max: "7.0.0", expected: true},
		{name: "7.0 above v6_x", version: "7.0.0", min: "6.1.0", max: "7.0.0", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ver, err := version.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("new version: %v", err)
			}

			min, err := version.NewVersion(tt.min)
			if err != nil {
				t.Fatalf("new min version: %v", err)
			}

			max, err := version.NewVersion(tt.max)
			if err != nil {
				t.Fatalf("new max version: %v", err)
			}

			if got := common.MatchesVersionRange(ver, min, max); got != tt.expected {
				t.Fatalf("MatchesVersionRange(%s, %s, %s) = %v, want %v", tt.version, tt.min, tt.max, got, tt.expected)
			}
		})
	}
}
