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
	Execute(ctx context.Context, filter regexfilter.RegexFilter, quite bool, latest bool) error
}

func NewGetScansCmd(cfg *config.Config, uc UseCaseGetScans) CmdGetScans {
	var (
		branchIdFlag string
		quite        bool
		latest       bool
		regexFilter  regexfilter.RegexFilter
	)

	cmd := &cobra.Command{
		Use:   "scans [<regex>]",
		Short: "Get AI scans",
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

			if err := uc.Execute(ctx, regexFilter, quite, latest); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'get scans' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "branch id")
	cmd.Flags().BoolVarP(&quite, "quite", "q", false, "Get only ids")
	cmd.Flags().BoolVar(&latest, "latest", false, "Get latest scan result only")

	return CmdGetScans{cmd}
}
