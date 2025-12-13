package create

import (
	"bytes"
	"testing"

	domainConfig "github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/create/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanCommand(t *testing.T) {
	t.Run("invalid UUID", func(t *testing.T) {
		cfg := &domainConfig.Config{}

		cmd := NewCreateBranchCommand(cfg, mocks.NewDependenciesContainer(t))
		buf := &bytes.Buffer{}
		cmd.SetOut(buf)
		cmd.SetErr(buf)

		args := []string{"-p", "not-a-uuid"}
		cmd.ParseFlags(args)

		err := cmd.Execute()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Validation error on field project-id: 'not-a-uuid' invalid uuid")
	})
}
