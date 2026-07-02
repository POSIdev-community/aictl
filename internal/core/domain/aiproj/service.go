package aiproj

import (
	"fmt"

	"github.com/POSIdev-community/aiproj/model"
	"github.com/POSIdev-community/aiproj/model/v1_10"
	"github.com/POSIdev-community/aiproj/model/v1_11"
	"github.com/POSIdev-community/aiproj/model/v1_9"
	"github.com/POSIdev-community/aiproj/versioning"
	"github.com/POSIdev-community/aiproj/versions"

	"github.com/POSIdev-community/aictl/internal/core/domain/version"
)

func ParseAndMigrateForServer(data []byte, serverVersion version.Version) (*Result, error) {
	target := getTargetVersion(serverVersion)

	engine, err := versioning.DefaultEngine()
	if err != nil {
		return nil, err
	}

	migrated, err := model.MigrateBytes(engine, data, target)
	if err != nil {
		return nil, fmt.Errorf("migrate aiproj to %s: %w", target, err)
	}

	switch target {
	case versions.V1_9:
		proj, err := v1_9.ReadFromBytes(migrated)
		if err != nil {
			return nil, fmt.Errorf("parse and migrate aiproj to %s: %w", target, err)
		}

		return &Result{Version: target, V19: &proj}, nil
	case versions.V1_10:
		proj, err := v1_10.ReadFromBytes(migrated)
		if err != nil {
			return nil, fmt.Errorf("parse and migrate aiproj to %s: %w", target, err)
		}

		return &Result{Version: target, V110: &proj}, nil
	case versions.V1_11:
		proj, err := v1_11.ReadFromBytes(migrated)
		if err != nil {
			return nil, fmt.Errorf("parse and migrate aiproj to %s: %w", target, err)
		}

		return &Result{Version: target, V111: &proj}, nil
	default:
		return nil, fmt.Errorf("unsupported aiproj version: %s", target)
	}
}
