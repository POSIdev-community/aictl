package context

import (
	"fmt"

	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"github.com/spf13/cobra"
)

func NewConfigShowCommand(
	cfg *config.Config,
	depsContainer *application.DependenciesContainer) *cobra.Command {

	var (
		json bool
		yaml bool
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current aictl context",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if json && yaml {
				return fmt.Errorf("cannot use both json and yaml flags")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			useCase, err := depsContainer.ConfigShowUseCase(cmd.Context())
			if err != nil {
				return err
			}

			err = useCase.Execute(cfg, json, yaml)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&json, "json", false, "Json format context")
	cmd.Flags().BoolVar(&yaml, "yaml", false, "Yaml format context")

	return cmd
}
