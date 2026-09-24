package get

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
	"github.com/POSIdev-community/aictl/internal/presenter/cmdtest"
)

type fakeReportTemplatesUC struct {
	called       int
	quiet        bool
	localization string
	filter       regexfilter.RegexFilter
}

func (f *fakeReportTemplatesUC) Execute(_ context.Context, filter regexfilter.RegexFilter, quiet bool, localization string) error {
	f.called++
	f.filter, f.quiet, f.localization = filter, quiet, localization
	return nil
}

func TestGetReportTemplatesCmd(t *testing.T) {
	t.Cleanup(resetGetPackageFlags)
	resetGetPackageFlags()
	cfg := mustGetCfg(t)
	uc := &fakeReportTemplatesUC{}
	ucs := defaultGetUCs()
	ucs.reportTemplates = uc
	root := buildGetRoot(t, cfg, ucs)
	require.NoError(t, cmdtest.Execute(t, root.Command, "report-templates", "owasp", "-q", "--localization", "ru"))
	require.True(t, uc.quiet)
	require.Equal(t, "ru", uc.localization)
}
