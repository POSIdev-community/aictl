package context

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewConfigSetCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		uriFlag       string
		tokenFlag     string
		tlsSkipFlag   bool
		noTlsSkipFlag bool
		projectIdFlag string
		branchIdFlag  string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set current aictl configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if uriFlag == "" && tokenFlag == "" && !tlsSkipFlag && !noTlsSkipFlag && projectIdFlag == "" && branchIdFlag == "" {
				return errs.NewValidationError("Any configs not provided")
			}

			if tlsSkipFlag && noTlsSkipFlag {
				return errs.NewValidationError("Not use both 'tls-skip' and 'no-tls-skip' flags at the same time")
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
				cfg.SetTLSSkip(true)
			}

			if noTlsSkipFlag {
				cfg.SetTLSSkip(false)
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
			log := logger.FromContext(cmd.Context())
			log.Info("aictl ctx")

			useCase, err := depsContainer.ConfigSetUseCase()
			if err != nil {
				return err
			}

			err = useCase.Execute(cfg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&uriFlag, "uri", "u", "", "AI server uri")
	cmd.Flags().StringVarP(&tokenFlag, "token", "t", "", "AI server access token")
	cmd.Flags().BoolVar(&tlsSkipFlag, "tls-skip", false, "Skip certificate verification")
	cmd.Flags().BoolVar(&noTlsSkipFlag, "no-tls-skip", false, "Not skip certificate verification")

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id")
	cmd.Flags().StringVarP(&branchIdFlag, "branch-id", "b", "", "Branch id")

	return cmd
}
