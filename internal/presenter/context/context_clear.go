package context

import (
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/spf13/cobra"
)

func NewConfigClearCommand(
	depsContainer *application.DependenciesContainer) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear current aictl configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			useCase, err := depsContainer.ConfigClearUseCase()
			if err != nil {
				return err
			}

			err = useCase.Execute()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
