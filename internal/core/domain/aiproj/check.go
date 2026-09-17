package aiproj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/POSIdev-community/aiproj/versioning"
	"github.com/POSIdev-community/aiproj/versions"
	"github.com/POSIdev-community/aiproj/versions/v1_10"
	"github.com/POSIdev-community/aiproj/versions/v1_11"
	"github.com/POSIdev-community/aiproj/versions/v1_8"
	"github.com/POSIdev-community/aiproj/versions/v1_9"
)

type Issue struct {
	Message string `json:"message"`
}

type FailureKind int

const (
	FailureNone FailureKind = iota
	FailureSchemaInvalid
	FailureVersionRequired
	FailureVersionMismatch
	FailureUnsupportedSchemaVersion
	FailureNotDetected
)

type Report struct {
	OK          bool        `json:"ok"`
	Version     string      `json:"version"`
	ProjectName *string     `json:"projectName"`
	Languages   *[]string   `json:"languages"`
	Errors      []Issue     `json:"errors"`
	Kind        FailureKind `json:"-"`
}

const (
	errLanguagesRequired = "(root): ProgrammingLanguages is required"
	errLanguagesNotArray = "(root): ProgrammingLanguages must be of type array"
	errLanguagesItemType = "(root): ProgrammingLanguages items must be of type string"
)

func SupportedSchemaVersions() []string {
	return []string{versions.V1_8, versions.V1_9, versions.V1_10, versions.V1_11}
}

func (k FailureKind) ShortMessage() string {
	switch k {
	case FailureVersionRequired:
		return "Version required"
	case FailureVersionMismatch:
		return "aiproj version mismatch"
	case FailureUnsupportedSchemaVersion:
		return "unsupported schema version"
	case FailureNotDetected:
		return "aiproj version was not detected"
	case FailureSchemaInvalid:
		return "aiproj schema invalid"
	default:
		return "aiproj check failed"
	}
}

func (r Report) HumanMessage() string {
	if r.OK {
		return ""
	}

	switch r.Kind {
	case FailureVersionRequired:
		return FailureVersionRequired.ShortMessage()
	case FailureVersionMismatch, FailureUnsupportedSchemaVersion, FailureNotDetected:
		if len(r.Errors) >= 1 {
			return r.Errors[0].Message
		}

		return r.Kind.ShortMessage()
	default:
		header := "aiproj schema is invalid:"
		if r.Version != "" {
			header = fmt.Sprintf("aiproj schema %s is invalid:", r.Version)
		}

		if len(r.Errors) == 0 {
			return header
		}

		var b strings.Builder
		b.WriteString(header)

		for _, issue := range r.Errors {
			b.WriteString("\n  - ")
			b.WriteString(issue.Message)
		}

		return b.String()
	}
}

// Check validates aiproj JSON. Empty schemaVersion means auto-detect via aiproj engine.
func Check(data []byte, schemaVersion string) (Report, error) {
	report, err := check(data, schemaVersion)
	if err != nil {
		return Report{}, err
	}

	return enrichReport(data, report), nil
}

func check(data []byte, schemaVersion string) (Report, error) {
	if schemaVersion != "" {
		if !isSupportedSchemaVersion(schemaVersion) {
			msg := fmt.Sprintf(
				"unsupported schema version %q (supported: %s)",
				schemaVersion,
				strings.Join(SupportedSchemaVersions(), ", "),
			)

			return Report{
				OK:      false,
				Version: "",
				Kind:    FailureUnsupportedSchemaVersion,
				Errors:  []Issue{{Message: msg}},
			}, nil
		}

		fileVersion := extractVersion(data)
		if fileVersion == "" {
			return Report{
				OK:      false,
				Version: "",
				Kind:    FailureVersionRequired,
				Errors:  []Issue{{Message: FailureVersionRequired.ShortMessage()}},
			}, nil
		}

		if fileVersion != schemaVersion {
			msg := fmt.Sprintf("aiproj Version is %q, expected %q", fileVersion, schemaVersion)

			return Report{
				OK:      false,
				Version: fileVersion,
				Kind:    FailureVersionMismatch,
				Errors:  []Issue{{Message: msg}},
			}, nil
		}

		if err := validateAs(data, schemaVersion); err != nil {
			return Report{
				OK:      false,
				Version: schemaVersion,
				Kind:    FailureSchemaInvalid,
				Errors:  ParseSchemaError(err),
			}, nil
		}

		return Report{OK: true, Version: schemaVersion, Errors: []Issue{}}, nil
	}

	engine, err := versioning.DefaultEngine()
	if err != nil {
		return Report{}, err
	}

	detected, err := engine.Detect(data)
	if err != nil {
		fileVersion := extractVersion(data)
		kind := FailureNotDetected
		if fileVersion != "" {
			kind = FailureSchemaInvalid
		}

		return Report{
			OK:      false,
			Version: fileVersion,
			Kind:    kind,
			Errors:  ParseSchemaError(err),
		}, nil
	}

	return Report{OK: true, Version: detected, Errors: []Issue{}}, nil
}

func enrichReport(data []byte, report Report) Report {
	report.ProjectName = extractProjectName(data)

	langs, langErr := extractLanguages(data)
	report.Languages = langs

	if langErr != "" {
		report.Errors = appendIssueIfAbsent(report.Errors, langErr)
		if report.OK {
			report.OK = false
			report.Kind = FailureSchemaInvalid
		}
	}

	return report
}

func extractProjectName(data []byte) *string {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	v, ok := raw["ProjectName"]
	if !ok || string(v) == "null" {
		return nil
	}

	var name string
	if err := json.Unmarshal(v, &name); err != nil {
		return nil
	}

	if name == "" {
		return nil
	}

	return &name
}

func extractLanguages(data []byte) (*[]string, string) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errLanguagesRequired
	}

	v, ok := raw["ProgrammingLanguages"]
	if !ok || string(v) == "null" {
		return nil, errLanguagesRequired
	}

	var items []json.RawMessage
	if err := json.Unmarshal(v, &items); err != nil {
		return nil, errLanguagesNotArray
	}

	langs := make([]string, 0, len(items))
	for _, item := range items {
		var s string
		if err := json.Unmarshal(item, &s); err != nil {
			return nil, errLanguagesItemType
		}

		langs = append(langs, s)
	}

	return &langs, ""
}

func appendIssueIfAbsent(issues []Issue, msg string) []Issue {
	for _, issue := range issues {
		if issue.Message == msg {
			return issues
		}
	}

	return append(issues, Issue{Message: msg})
}

func isSupportedSchemaVersion(v string) bool {
	for _, s := range SupportedSchemaVersions() {
		if s == v {
			return true
		}
	}

	return false
}

func extractVersion(data []byte) string {
	var base struct {
		Version string `json:"Version"`
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&base); err != nil {
		return ""
	}

	return base.Version
}

func validateAs(data []byte, ver string) error {
	switch ver {
	case versions.V1_8:
		return v1_8.Validate(data)
	case versions.V1_9:
		return v1_9.Validate(data)
	case versions.V1_10:
		return v1_10.Validate(data)
	case versions.V1_11:
		return v1_11.Validate(data)
	default:
		return fmt.Errorf("unsupported schema version: %s", ver)
	}
}

// ParseSchemaError extracts schema issue messages from aiproj library error strings.
func ParseSchemaError(err error) []Issue {
	if err == nil {
		return nil
	}

	text := err.Error()
	const marker = "schema validation failed: "

	if idx := strings.LastIndex(text, marker); idx >= 0 {
		text = text[idx+len(marker):]
	}

	parts := strings.Split(text, "; ")
	issues := make([]Issue, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		issues = append(issues, Issue{Message: part})
	}

	if len(issues) == 0 {
		return []Issue{{Message: err.Error()}}
	}

	return issues
}
