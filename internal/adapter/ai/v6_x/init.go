package v6_x

import (
	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

var (
	minVersion, _ = version.NewVersion("6.3.0")
	maxVersion, _ = version.NewVersion("7.0.0")
)

var Initializer common.Initializer = common.VersionRangeInitializer{
	MinVersion: minVersion,
	MaxVersion: maxVersion,
	NewClient: func(base *common.BaseClient) common.ClientAi {
		return NewAiClient(base)
	},
}
