package application

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/POSIdev-community/aictl/internal/adapter/config"
	"github.com/POSIdev-community/aictl/internal/core/application"
	"github.com/POSIdev-community/aictl/internal/presenter"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type Application struct {
	cmd *cobra.Command
}

func NewApplication() *Application {
	cfgAdapter := config.NewContextAdapter()
	cfg := cfgAdapter.GetContextFromAictlFolder()

	dependencyContainer := application.NewDependenciesContainer(cfgAdapter)

	cmd := presenter.NewRootCmd(cfg, dependencyContainer)
	cmd.DisableAutoGenTag = true

	return &Application{cmd}
}

func (app *Application) Run(ctx context.Context) {
	err := app.cmd.ExecuteContext(ctx)
	if err == nil {
		os.Exit(0)
	}
	log := logger.FromContext(ctx)
	log.StdErrF(err.Error())

	exitCode, errorMessage := errs.MapExitCode(err)

	_, err = fmt.Fprintln(os.Stderr, errorMessage)

	os.Exit(exitCode)
}

func (app *Application) GenerateDoc(path string) error {
	if err := doc.GenMarkdownTree(app.cmd, path); err != nil {
		return fmt.Errorf("err while generate doc: %w", err)
	}

	return nil
}
