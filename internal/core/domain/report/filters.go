package report

import (
	"fmt"
	"slices"

	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

// Filters holds optional vulnerability filters for report generation.
// When Apply is false, the API is called with useFilters=false and filters omitted.
// When Apply is true, useFilters=true; only non-nil bool pointers are sent as true;
// Types, Languages, and ScanModules are always sent (possibly empty slices).
type Filters struct {
	Apply bool

	LevelHigh      *bool
	LevelMedium    *bool
	LevelLow       *bool
	LevelPotential *bool

	StatusUndefined     *bool
	StatusConfirmed     *bool
	StatusConfirmedAuto *bool
	StatusRejected      *bool

	ModeEntryPoint    *bool
	ModePublicMethods *bool
	ModeRootFunction  *bool
	ModeOthers        *bool

	FoundThisScan *bool
	FoundPrevScan *bool

	Conditional    *bool
	NonConditional *bool
	Suppressed     *bool
	NonSuppressed  *bool

	Suspected    *bool
	SecondLevel  *bool
	OnlyFavorite *bool

	Types       []string
	Languages   []string
	ScanModules []string
}

// AllowedScanModules is the CLI whitelist for --scan-module (exact API case).
var AllowedScanModules = []string{
	"StaticCodeAnalysis",
	"PatternMatching",
	"Components",
	"SoftwareCompositionAnalysis",
	"Configuration",
	"MaliciousCodeDetection",
	"SecretDetection",
	"BlackBox",
}

// AllowedLanguages is the CLI whitelist for --language (exact API case, without None).
var AllowedLanguages = []string{
	"CAndCPlusPlus",
	"CSharp",
	"Dart",
	"Go",
	"Java",
	"JavaScript",
	"Kotlin",
	"ObjectiveC",
	"OneC",
	"Php",
	"Python",
	"Ruby",
	"Scala",
	"Solidity",
	"Sql",
	"Swift",
}

// EmptyFilters returns filters that disable useFilters on the API.
func EmptyFilters() Filters {
	return Filters{}
}

// HasAny reports whether any filter field is set (for CLI “≥1 flag required”).
func (f Filters) HasAny() bool {
	if f.LevelHigh != nil || f.LevelMedium != nil || f.LevelLow != nil || f.LevelPotential != nil {
		return true
	}
	if f.StatusUndefined != nil || f.StatusConfirmed != nil || f.StatusConfirmedAuto != nil || f.StatusRejected != nil {
		return true
	}
	if f.ModeEntryPoint != nil || f.ModePublicMethods != nil || f.ModeRootFunction != nil || f.ModeOthers != nil {
		return true
	}
	if f.FoundThisScan != nil || f.FoundPrevScan != nil {
		return true
	}
	if f.Conditional != nil || f.NonConditional != nil || f.Suppressed != nil || f.NonSuppressed != nil {
		return true
	}
	if f.Suspected != nil || f.SecondLevel != nil || f.OnlyFavorite != nil {
		return true
	}
	return len(f.Types) > 0 || len(f.Languages) > 0 || len(f.ScanModules) > 0
}

// Normalize deduplicates array fields preserving first-seen order.
func (f *Filters) Normalize() {
	f.Types = dedupePreserveOrder(f.Types)
	f.Languages = dedupePreserveOrder(f.Languages)
	f.ScanModules = dedupePreserveOrder(f.ScanModules)
}

// ValidateArrays checks type/language/scan-module values after CLI parsing.
func (f Filters) ValidateArrays() error {
	for _, t := range f.Types {
		if t == "" {
			return validation.NewFieldError("type", "cannot be empty")
		}
	}
	for _, lang := range f.Languages {
		if lang == "None" {
			return validation.NewFieldError("language", "None is not allowed")
		}
		if !slices.Contains(AllowedLanguages, lang) {
			return validation.NewFieldError("language", fmt.Sprintf("unsupported value %q", lang))
		}
	}
	for _, m := range f.ScanModules {
		if !slices.Contains(AllowedScanModules, m) {
			return validation.NewFieldError("scan-module", fmt.Sprintf("unsupported value %q", m))
		}
	}
	return nil
}

func dedupePreserveOrder(in []string) []string {
	if len(in) == 0 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
