package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	domainlicense "github.com/POSIdev-community/aictl/internal/core/domain/license"
)

func TestLicenseFromAPI_LegacyVsNew(t *testing.T) {
	t.Parallel()

	legacy := common.LicenseFromAPI([]string{"Go"}, nil, false)
	require.Nil(t, legacy.LicensedModules)
	require.ElementsMatch(t, domainlicense.LicensedScanModuleTypes(), legacy.AllowedScanModules())

	emptyNew := common.LicenseFromAPI([]string{"Go"}, nil, true)
	require.NotNil(t, emptyNew.LicensedModules)
	require.Empty(t, emptyNew.AllowedScanModules())

	withModules := common.LicenseFromAPI([]string{"Go"}, []common.LicensedModuleInput{
		{ID: "Sca", Enabled: true, ScanModuleTypes: []string{domainlicense.ScanModuleSoftwareCompositionAnalysis}},
	}, true)
	require.Equal(t, []string{domainlicense.ScanModuleSoftwareCompositionAnalysis}, withModules.AllowedScanModules())
}

func TestMapLanguageGroups(t *testing.T) {
	t.Parallel()

	require.Nil(t, common.MapLanguageGroups[string](nil))
	langs := []string{"Go", "Java"}
	require.Equal(t, []string{"Go", "Java"}, common.MapLanguageGroups(&langs))
}
