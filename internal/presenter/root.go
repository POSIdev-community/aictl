package presenter

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
	"github.com/POSIdev-community/aictl/internal/presenter/context"
	"github.com/POSIdev-community/aictl/internal/presenter/create"
	deletePresenter "github.com/POSIdev-community/aictl/internal/presenter/delete"
	"github.com/POSIdev-community/aictl/internal/presenter/get"
	"github.com/POSIdev-community/aictl/internal/presenter/scan"
	"github.com/POSIdev-community/aictl/internal/presenter/set"
	"github.com/POSIdev-community/aictl/internal/presenter/update"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/POSIdev-community/aictl/pkg/version"
)

type CmdRoot struct {
	*cobra.Command
}

func NewRootCmd(contextCmd *context.CmdContext, createCmd *create.CmdCreate, deleteCmd *deletePresenter.CmdDelete,
	getCmd *get.CmdGet, scanCmd *scan.CmdScan, setCmd *set.CmdSet, updateCmd *update.CmdUpdate) *CmdRoot {

	var versionFlag bool

	rootCmd := &cobra.Command{
		Use:   "aictl",
		Short: "Application Inspector ConTroL tool",
		Long: `CLI for managing PT Application Inspector: projects, branches, scans, reports, and local context.

Use subcommands to talk to an AI server. Connection settings come from context (~/.config/aictl/context.yaml) or -u/-t/--tls-skip flags.`,
		Example: `  aictl ctx set -u https://ai.example -t <token>
  aictl get projects
  aictl --version`,
		PersistentPreRunE: _utils.InitializeLogger,
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionFlag {
				l := logger.FromContext(cmd.Context())
				l.StdOut(version.GetVersion())

				return nil
			}

			if len(args) == 0 {
				err := cmd.Help()
				if err != nil {
					return fmt.Errorf("get help: %w", err)
				}
				return nil
			}

			return nil
		},
	}

	rootCmd.AddCommand(contextCmd.Command)
	rootCmd.AddCommand(createCmd.Command)
	rootCmd.AddCommand(deleteCmd.Command)
	rootCmd.AddCommand(getCmd.Command)
	rootCmd.AddCommand(scanCmd.Command)
	rootCmd.AddCommand(setCmd.Command)
	rootCmd.AddCommand(updateCmd.Command)

	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "Show aictl version")

	return &CmdRoot{rootCmd}
}
