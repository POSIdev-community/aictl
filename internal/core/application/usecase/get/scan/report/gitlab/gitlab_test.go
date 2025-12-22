package gitlab_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	. "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/gitlab"
	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/sarif/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
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

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
		aiAdapter.On("GetTemplateId", t.Context(), report.GitlabReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)
		cliAdapter.On("ShowReader", reportReader).Return(nil).Once()
		cliAdapter.On("ShowTextf", t.Context(), "getting gitlab scan report, id '%v'", []interface{}{scanID.String()}).Return().Once()
		cliAdapter.On("ShowText", t.Context(), "gitlab scan report got").Return().Once()

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), scanID, "", includeComments, includeDfd, includeGlossary))
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

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("InitializeWithRetry", t.Context()).Return(nil).Once()
		aiAdapter.On("GetTemplateId", t.Context(), report.GitlabReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)
		cliAdapter.On("ShowTextf", t.Context(), "getting gitlab scan report, id '%v'", []interface{}{scanID.String()}).Return().Once()
		cliAdapter.On("ShowText", t.Context(), "gitlab scan report got").Return().Once()

		cfg := config.NewConfig(config.Uri{}, "", true, projectID, uuid.New())

		uc, err := NewUseCase(aiAdapter, cliAdapter, cfg)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), scanID, filePath, includeComments, includeDfd, includeGlossary))

		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, reportText, string(data))
	})
}
