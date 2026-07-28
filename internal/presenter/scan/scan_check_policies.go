package scan

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
	"github.com/POSIdev-community/aictl/internal/presenter/.utils"
)

type CmdScanCheckPolicies struct {
	*cobra.Command
}

type UseCaseScanCheckPolicies interface {
	Execute(ctx context.Context, scanId uuid.UUID, failOnPoliciesRejected bool) error
}

func NewScanCheckPoliciesCmd(cfg *config.Config, uc UseCaseScanCheckPolicies) CmdScanCheckPolicies {
	var (
		projectIdFlag          string
		scanIdFlag             string
		scanId                 uuid.UUID
		failOnPoliciesRejected bool
	)

	cmd := &cobra.Command{
		Use:   "check-policies <scan-id>",
		Short: "Check scan policy state",
		Long:  `Check the policy state for a scan. Project id comes from context or -p.`,
		Example: `  aictl scan check-policies <scan-id>
  aictl scan check-policies <scan-id> --fail-on-policies-rejected`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = cfg.UpdateProjectId(projectIdFlag); err != nil {
				return err
			}

			args = _utils.ReadArgsFromStdin(args)
			if len(args) < 1 {
				return validation.NewError("missing scan id")
			}

			scanIdFlag = args[0]
			scanId, err = uuid.Parse(scanIdFlag)
			if err != nil {
				return validation.NewFieldError(scanIdFlag, "invalid uuid")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := uc.Execute(ctx, scanId, failOnPoliciesRejected); err != nil {
				cmd.SilenceUsage = true

				return fmt.Errorf("'scan check-policies' usecase call: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectIdFlag, "project-id", "p", "", "Project id (overrides context)")
	cmd.Flags().BoolVar(&failOnPoliciesRejected, "fail-on-policies-rejected", false, "Exit 1 if PolicyState is Rejected")

	return CmdScanCheckPolicies{cmd}
}
