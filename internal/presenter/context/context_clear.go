package context

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewConfigClearCommand(
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear current aictl configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.FromContext(cmd.Context())
			log.StdErr("start clearing context")

			useCase, err := depsContainer.ConfigClearUseCase(cmd.Context())
			if err != nil {
				return err
			}

			err = useCase.Execute(skipConfirm)
			if err != nil {
				return err
			}

			log.StdErr("context successfully cleared")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
