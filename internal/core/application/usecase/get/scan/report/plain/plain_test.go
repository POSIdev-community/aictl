package plain_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	. "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/plain"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/plain/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("stdout write", func(t *testing.T) {
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

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
		aiAdapter.On("GetTemplateId", t.Context(), report.PlainReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary, l10n).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)
		cliAdapter.On("ShowReader", reportReader).Return(nil).Once()
		cliAdapter.On("ShowTextf", t.Context(), "getting plain scan report, id '%v'", []interface{}{scanID.String()}).Return().Once()
		cliAdapter.On("ShowText", t.Context(), "plain scan report got").Return().Once()

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), scanID, "", includeComments, includeDfd, includeGlossary, l10n))
	})

	t.Run("write to file", func(t *testing.T) {
		t.Parallel()

		projectID := uuid.New()
		scanID := uuid.New()
		templateID := uuid.New()
		reportText := "foo: BAR"
		reportReader := io.NopCloser(bytes.NewBufferString(reportText))
		filePath := filepath.Join(t.TempDir(), "test.txt")
		includeComments := false
		includeDfd := false
		includeGlossary := false
		l10n := "en"

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
		aiAdapter.On("GetTemplateId", t.Context(), report.PlainReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary, l10n).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)
		cliAdapter.On("ShowTextf", t.Context(), "getting plain scan report, id '%v'", []interface{}{scanID.String()}).Return().Once()
		cliAdapter.On("ShowText", t.Context(), "plain scan report got").Return().Once()

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), scanID, filePath, includeComments, includeDfd, includeGlossary, l10n))

		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, reportText, string(data))
	})
}
