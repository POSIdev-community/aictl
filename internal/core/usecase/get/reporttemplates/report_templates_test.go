package reporttemplates

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
)

func TestUseCase_Execute_FilterAndQuite(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	templates := []report.Template{
		{Id: uuid.New(), Name: "SARIF"},
		{Id: uuid.New(), Name: "HTML"},
	}

	ai := NewMockAI(t)
	ai.On("InitializeWithRetry", ctx).Return(nil).Once()
	ai.On("GetReportTemplates", ctx, "en").Return(templates, nil).Once()

	filter, err := regexfilter.NewRegexFilter("SARIF")
	require.NoError(t, err)

	cli := NewMockCLI(t)
	cli.On("ShowReportTemplatesQuite", ctx, []report.Template{templates[0]}).Return().Once()

	uc, err := NewUseCase(ai, cli)
	require.NoError(t, err)
	require.NoError(t, uc.Execute(ctx, filter, true, "en"))
}
