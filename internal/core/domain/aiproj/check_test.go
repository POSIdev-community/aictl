package aiproj

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aiproj/versions/v1_11"
)

func TestCheck_ValidFixture(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "")
	require.NoError(t, err)
	require.True(t, report.OK)
	require.Equal(t, "1.11", report.Version)
	require.Empty(t, report.Errors)
	require.NotNil(t, report.ProjectName)
	require.Equal(t, "demo", *report.ProjectName)
	require.NotNil(t, report.Languages)
	require.Equal(t, []string{"Go"}, *report.Languages)
}

func TestCheck_SchemaInvalid(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": 123,
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.Equal(t, FailureSchemaInvalid, report.Kind)
	require.Equal(t, "1.11", report.Version)
	require.NotEmpty(t, report.Errors)
	require.Nil(t, report.ProjectName)
	require.NotNil(t, report.Languages)
	require.Equal(t, []string{"Go"}, *report.Languages)
	require.Contains(t, report.HumanMessage(), "aiproj schema 1.11 is invalid:")
	require.Equal(t, "aiproj schema invalid", report.Kind.ShortMessage())
}

func TestCheck_VersionRequired(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.Equal(t, FailureVersionRequired, report.Kind)
	require.Equal(t, "Version required", report.HumanMessage())
	require.NotNil(t, report.ProjectName)
	require.Equal(t, "demo", *report.ProjectName)
	require.NotNil(t, report.Languages)
	require.Equal(t, []string{"Go"}, *report.Languages)
}

func TestCheck_VersionMismatch(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.10",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go"],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.Equal(t, FailureVersionMismatch, report.Kind)
	require.Equal(t, `aiproj Version is "1.10", expected "1.11"`, report.HumanMessage())
	require.Len(t, report.Errors, 1)
	require.NotNil(t, report.ProjectName)
	require.Equal(t, "demo", *report.ProjectName)
	require.NotNil(t, report.Languages)
	require.Equal(t, []string{"Go"}, *report.Languages)
}

func TestCheck_UnsupportedSchemaVersion(t *testing.T) {
	t.Parallel()

	report, err := Check([]byte(`{"Version":"1.11"}`), "1.12")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.Equal(t, FailureUnsupportedSchemaVersion, report.Kind)
	require.Contains(t, report.HumanMessage(), `unsupported schema version "1.12"`)
	require.Nil(t, report.ProjectName)
	require.Nil(t, report.Languages)
	require.Contains(t, issueMessages(report), errLanguagesRequired)
}

func TestCheck_EmptyProjectNameAndLanguages(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": "",
		"ProgrammingLanguages": [],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.Nil(t, report.ProjectName)
	require.NotNil(t, report.Languages)
	require.Empty(t, *report.Languages)
}

func TestCheck_LanguagesMixedTypes(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": ["Go", 1],
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.NotNil(t, report.ProjectName)
	require.Equal(t, "demo", *report.ProjectName)
	require.Nil(t, report.Languages)
	require.Contains(t, issueMessages(report), errLanguagesItemType)
}

func TestCheck_LanguagesNotArray(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": "demo",
		"ProgrammingLanguages": "Go",
		"ScanModules": ["StaticCodeAnalysis"]
	}`)

	report, err := Check(data, "1.11")
	require.NoError(t, err)
	require.False(t, report.OK)
	require.Nil(t, report.Languages)
	require.Contains(t, issueMessages(report), errLanguagesNotArray)
}

func issueMessages(report Report) []string {
	msgs := make([]string, 0, len(report.Errors))
	for _, issue := range report.Errors {
		msgs = append(msgs, issue.Message)
	}

	return msgs
}

func TestParseSchemaError(t *testing.T) {
	t.Parallel()

	err := errors.New("validate aiproj 1.11 schema: schema validation failed: (root): ProjectName must be of type string; ScanModules.0: must be one of the following")
	issues := ParseSchemaError(err)
	require.Len(t, issues, 2)
	require.Equal(t, "(root): ProjectName must be of type string", issues[0].Message)
	require.Equal(t, "ScanModules.0: must be one of the following", issues[1].Message)
}

func TestParseSchemaError_FromRealValidate(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"Version": "1.11",
		"ProjectName": 123,
		"ProgrammingLanguages": ["NotALang"],
		"ScanModules": ["Nope"]
	}`)

	err := v1_11.Validate(data)
	require.Error(t, err)

	issues := ParseSchemaError(err)
	require.NotEmpty(t, issues)
}
