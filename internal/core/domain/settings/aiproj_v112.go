package settings

import (
	"encoding/json"
	"fmt"

	v111 "github.com/POSIdev-community/aiproj/model/v1_11"
	v112 "github.com/POSIdev-community/aiproj/model/v1_12"
)

func (s *ScanSettings) UpdateFromV112(p *v112.AIProj) error {
	if p == nil {
		return fmt.Errorf("aiproj 1.12 payload is nil")
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal aiproj 1.12: %w", err)
	}

	var converted v111.AIProj
	if err := json.Unmarshal(data, &converted); err != nil {
		return fmt.Errorf("convert aiproj 1.12 to 1.11 model: %w", err)
	}

	if err := s.UpdateFromV111(&converted); err != nil {
		return err
	}

	if p.AutomatedAiTriage != nil {
		s.AutomatedAiTriageEnabled = *p.AutomatedAiTriage
	}

	if p.AutomatedAiTriageMinSeverity != nil {
		s.AutomatedAiTriageMinSeverity = string(*p.AutomatedAiTriageMinSeverity)
	}

	return nil
}
