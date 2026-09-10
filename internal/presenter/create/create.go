package create

import (
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdCreate struct {
	*cobra.Command
}

var safeFlag bool

func NewCreateCmd(
	cfg *config.Config,
	cmdCreateBranch CmdCreateBranch,
	cmdCreateProject CmdCreateProject,
	cmdCreateSbomProject CmdCreateSbomProject) *CmdCreate {

	cmd := &cobra.Command{
		Use:               "create",
		Short:             "Create resources",
		Long:              `Create AI projects (source or SBOM) and branches on the server.`,
		PersistentPreRunE: _utils.ChainRunE(_utils.InitializeLogger, _utils.UpdateConfig(cfg)),
	}

	cmd.AddCommand(cmdCreateProject.Command)
	cmd.AddCommand(cmdCreateSbomProject.Command)
	cmd.AddCommand(cmdCreateBranch.Command)

	_utils.AddConnectionPersistentFlags(cmd)

	cmd.PersistentFlags().BoolVar(&safeFlag, "safe", false, "If the resource already exists, return its id without error")

	return &CmdCreate{cmd}
}
