package settings

import (
	"encoding/json"
	"fmt"

	v110 "github.com/POSIdev-community/aiproj/model/v1_10"
	v111 "github.com/POSIdev-community/aiproj/model/v1_11"
)

func (s *ScanSettings) UpdateFromV111(p *v111.AIProj) error {
	if p == nil {
		return fmt.Errorf("aiproj 1.11 payload is nil")
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal aiproj 1.11: %w", err)
	}

	var converted v110.AIProj
	if err := json.Unmarshal(data, &converted); err != nil {
		return fmt.Errorf("convert aiproj 1.11 to 1.10 model: %w", err)
	}

	return s.UpdateFromV110(&converted)
}
