package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

type CmdConfigSet struct {
	*cobra.Command
}

type UseCaseConfigSet interface {
	Execute() error
}

func NewConfigSetCommand(cfg *config.Config, uc UseCaseConfigSet) CmdConfigSet {

	var (
		uriFlag       string
		tokenFlag     string
		tlsSkipFlag   bool
		noTlsSkipFlag bool
		cacertFlag    string
		projectIdFlag string
		branchIdFlag  string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set current aictl configuration",
		Long:  `Update one or more fields in the local aictl context. At least one flag is required. Do not pass both --tls-skip and --no-tls-skip. Do not pass --cacert together with --tls-skip. Clear cacert with: aictl ctx unset --cacert.`,
		Example: `  aictl ctx set -u https://ai.example -t <token>
  aictl ctx set -p <project-id> -b <branch-id>
  aictl ctx set --cacert /path/to/ca.pem
  aictl ctx set --tls-skip`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if uriFlag == "" && tokenFlag == "" && !tlsSkipFlag && !noTlsSkipFlag && cacertFlag == "" && projectIdFlag == "" && branchIdFlag == "" {
				return validation.NewError("Any configs not provided")
			}

			if tlsSkipFlag && noTlsSkipFlag {
				return validation.NewError("Not use both 'tls-skip' and 'no-tls-skip' flags at the same time")
			}

			if tlsSkipFlag && cacertFlag != "" {
				return validation.NewError("Not use both 'cacert' and 'tls-skip' flags at the same time")
			}

			if uriFlag != "" {
				if err := cfg.SetURI(uriFlag); err != nil {
					return err
				}
			}

			if tokenFlag != "" {
				if err := cfg.SetToken(tokenFlag); err != nil {
					return err
				}
			}

			if tlsSkipFlag {
				if cfg.CACertPath() != "" {
					return validation.NewError("cannot enable 'tls-skip' while 'cacert' is set; unset cacert first")
				}
				cfg.SetTLSSkip(true)
			}

			if noTlsSkipFlag {
				cfg.SetTLSSkip(false)
			}

			if cacertFlag != "" {
				if _, err := os.Stat(cacertFlag); err != nil {
					return validation.NewFieldError("cacert", fmt.Sprintf("file %q: %v", cacertFlag, err))
				}
				if err := cfg.SetCACertPath(cacertFlag); err != nil {
					return err
				}
			}

			if projectIdFlag != "" {
				if err := cfg.SetProjectId(projectIdFlag); err != nil {
					return err
				}
			}

			if branchIdFlag != "" {
				if err := cfg.SetBranchId(branchIdFlag); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uc.Execute()
			if err != nil {
				return fmt.Errorf("'ctx set' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&uriFlag, "uri", "u", "", "AI server URI")
	cmd.Flags().StringVarP(&tokenFlag, "token", "t", "", "AI server access token")
	cmd.Flags().BoolVar(&tlsSkipFlag, "tls-skip", false, "Skip TLS certificate verification")
	cmd.Flags().BoolVar(&noTlsSkipFlag, "no-tls-skip", false, "Require TLS certificate verification")
	cmd.Flags().StringVar(&cacertFlag, "cacert", "", "Path to PEM file with CA certificate(s) to trust (appended to system roots)")

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Default project id")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "Default branch id")

	return CmdConfigSet{cmd}
}
