package sarif_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/sarif/mocks"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/port"

	. "github.com/POSIdev-community/aictl/internal/core/application/usecase/get/scan/report/sarif"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	t.Run("stdout write", func(t *testing.T) {
		t.Parallel()

		projectID := uuid.New()
		scanID := uuid.New()
		templateID := uuid.New()
		report := "foo: BAR"
		reportReader := io.NopCloser(bytes.NewBufferString(report))
		includeComments := false
		includeDfd := false
		includeGlossary := false

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("GetTemplateId", t.Context(), port.SarifReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)
		cliAdapter.On("ShowReader", reportReader).Return(nil).Once()

		uc, err := NewUseCase(aiAdapter, cliAdapter)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), config.NewConfig(config.Uri{}, "", true, projectID, uuid.New()), scanID, "", includeComments, includeDfd, includeGlossary))
	})

	t.Run("write to file", func(t *testing.T) {
		t.Parallel()

		projectID := uuid.New()
		scanID := uuid.New()
		templateID := uuid.New()
		report := "foo: BAR"
		reportReader := io.NopCloser(bytes.NewBufferString(report))
		filePath := filepath.Join(t.TempDir(), "test.txt")
		includeComments := false
		includeDfd := false
		includeGlossary := false

		aiAdapter := mocks.NewAI(t)
		aiAdapter.On("GetTemplateId", t.Context(), port.SarifReportType).Return(templateID, nil).Once()
		aiAdapter.On("GetReport", t.Context(), projectID, scanID, templateID, includeComments, includeDfd, includeGlossary).Return(reportReader, nil).Once()

		cliAdapter := mocks.NewCLI(t)

		uc, err := NewUseCase(aiAdapter, cliAdapter)
		require.NoError(t, err)

		require.NoError(t, uc.Execute(t.Context(), config.NewConfig(config.Uri{}, "", true, projectID, uuid.New()), scanID, filePath, includeComments, includeDfd, includeGlossary))

		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, report, string(data))
	})
}
