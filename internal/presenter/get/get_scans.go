package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/regexfilter"
)

type CmdGetScans struct {
	*cobra.Command
}

type UseCaseGetScans interface {
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quiet bool, latest bool) error
}

func NewGetScansCmd(cfg *config.Config, uc UseCaseGetScans, cmdSbomProject CmdGetScansSbomProject) CmdGetScans {
	var (
		branchIdFlag string
		quiet        bool
		latest       bool
		regexFilter  regexfilter.RegexFilter
	)

	cmd := &cobra.Command{
		Use:   "scans [<regex>]",
		Short: "Get AI scans",
		Long:  `List scans for the current branch matching an optional regex filter. Branch id comes from context or -b.`,
		Example: `  aictl get scans -b <branch-id>
  aictl get scans 2024 -b <branch-id>
  aictl get scans --latest -b <branch-id>`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.UpdateBranchId(branchIdFlag); err != nil {
				return err
			}

			filter := strings.Join(args, " ")

			var err error
			regexFilter, err = regexfilter.NewRegexFilter(filter)
			if err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("new filter: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, regexFilter, quiet, latest); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scans' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "Branch id (overrides context)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Print only ids")
	cmd.Flags().BoolVar(&latest, "latest", false, "Return only the latest scan result")

	cmd.AddCommand(cmdSbomProject.Command)

	return CmdGetScans{cmd}
}
