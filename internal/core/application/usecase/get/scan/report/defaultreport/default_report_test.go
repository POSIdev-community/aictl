package default_report_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	defaultreport "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/defaultreport"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/defaultreport/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	reportTypes := []report.ReportType{
		report.AutoCheck,
		report.Custom,
		report.Gitlab,
		report.Json,
		report.Markdown,
		report.Nist,
		report.Oud4,
		report.Owasp,
		report.Owaspm,
		report.Pcidss,
		report.PlainReport,
		report.Sans,
		report.Sarif,
		report.Xml,
	}

	t.Run("stdout write", func(t *testing.T) {
		t.Parallel()

		for _, rt := range reportTypes {
			t.Run(rt.String(), func(t *testing.T) {
				t.Parallel()

				projectID := uuid.New()
				scanID := uuid.New()
				templateID := uuid.New()
				reportText := "foo: BAR"
				reportReader := io.NopCloser(bytes.NewBufferString(reportText))
				includeComments := false
				includeDfd := false
				includeGlossary := false
				l10n := "en"
				reportType := rt

				aiAdapter := mocks.NewAI(t)
				aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
				aiAdapter.On("GetDefaultTemplateId", t.Context(), reportType).Return(templateID, nil).Once()
				aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary, l10n).Return(reportReader, nil).Once()

				cliAdapter := mocks.NewCLI(t)
				cliAdapter.On("ShowReader", reportReader).Return(nil).Once()
				cliAdapter.On("ShowTextf", t.Context(), "getting '%s' scan report, scan-id '%v'", []interface{}{reportType.String(), scanID.String()}).Return().Once()
				cliAdapter.On("ShowTextf", t.Context(), "'%s' scan report got", []interface{}{reportType.String()}).Return().Once()

				cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

				uc, err := defaultreport.NewUseCase(aiAdapter, cliAdapter, cfg)
				require.NoError(t, err)

				require.NoError(t, uc.Execute(t.Context(), scanID, reportType, "", includeComments, includeDfd, includeGlossary, l10n))
			})
		}
	})

	t.Run("write to file", func(t *testing.T) {
		t.Parallel()

		for _, rt := range reportTypes {
			t.Run(rt.String(), func(t *testing.T) {
				t.Parallel()

				projectID := uuid.New()
				scanID := uuid.New()
				templateID := uuid.New()
				reportText := "foo: BAR"
				reportReader := io.NopCloser(bytes.NewBufferString(reportText))
				filePath := filepath.Join(t.TempDir(), "test_"+rt.String()+".txt")
				includeComments := false
				includeDfd := false
				includeGlossary := false
				l10n := "en"
				reportType := rt

				aiAdapter := mocks.NewAI(t)
				aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
				aiAdapter.On("GetDefaultTemplateId", t.Context(), reportType).Return(templateID, nil).Once()
				aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary, l10n).Return(reportReader, nil).Once()

				cliAdapter := mocks.NewCLI(t)
				cliAdapter.On("ShowTextf", t.Context(), "getting '%s' scan report, scan-id '%v'", []interface{}{reportType.String(), scanID.String()}).Return().Once()
				cliAdapter.On("ShowTextf", t.Context(), "'%s' scan report got", []interface{}{reportType.String()}).Return().Once()

				cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

				uc, err := defaultreport.NewUseCase(aiAdapter, cliAdapter, cfg)
				require.NoError(t, err)

				require.NoError(t, uc.Execute(t.Context(), scanID, reportType, filePath, includeComments, includeDfd, includeGlossary, l10n))

				data, err := os.ReadFile(filePath)
				require.NoError(t, err)
				require.Equal(t, reportText, string(data))
			})
		}
	})
}
