package aiproj

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/POSIdev-community/aictl/internal/core/domain/aiproj/internal/1_8"
)

type AIProj = aiproj1_8.AIProj

// FromString parses a JSON string into an AIProj structure and returns it, or an error if decoding fails.
func FromString(s string) (AIProj, error) {
	reader := bytes.NewBufferString(s)

	return FromReader(reader)
}

// FromReader reads JSON data from the provided io.Reader and decodes it into an AIProj structure.
// It returns the AIProj instance or an error if decoding fails or unknown fields are present in the input.
func FromReader(r io.Reader) (AIProj, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var aiProj AIProj
	if err := decoder.Decode(&aiProj); err != nil {
		return AIProj{}, err
	}

	return aiProj, nil
}
