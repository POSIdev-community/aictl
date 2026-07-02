package aiproj

import (
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
	"github.com/POSIdev-community/aiproj/versions"
)

var (
	targetV19, _  = version.NewVersion(versions.V1_9)
	targetV110, _ = version.NewVersion(versions.V1_10)
	targetV111, _ = version.NewVersion(versions.V1_11)

	serverV60, _ = version.NewVersion("6.0.0")
	serverV61, _ = version.NewVersion("6.1.0")
)

func getTargetVersion(serverVersion version.Version) string {
	switch {
	case !serverVersion.Less(serverV61):
		return targetV111.String()
	case !serverVersion.Less(serverV60):
		return targetV110.String()
	default:
		return targetV19.String()
	}
}
