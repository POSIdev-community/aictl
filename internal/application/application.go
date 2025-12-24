package application

import (
	"context"
	"fmt"
	"os"

	"github.com/POSIdev-community/aictl/internal/di"
	"github.com/POSIdev-community/aictl/internal/presenter"
	"github.com/spf13/cobra/doc"

	"github.com/POSIdev-community/aictl/internal/adapter/config"
	"github.com/POSIdev-community/aictl/pkg/errs"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type Application struct {
	cmd *presenter.CmdRoot
}

func NewApplication() *Application {
	cfgAdapter := config.NewContextAdapter()
	cfg := cfgAdapter.GetContextFromAictlFolder()

	cmd, _ := di.InitializeCmd(cfg)
	cmd.DisableAutoGenTag = true

	return &Application{cmd}
}

func (app *Application) Run(ctx context.Context) {
	err := app.cmd.ExecuteContext(ctx)
	if err == nil {
		os.Exit(0)
	}
	log := logger.FromContext(ctx)
	log.StdErrf(err.Error())

	exitCode, errorMessage := errs.MapExitCode(err)

	_, err = fmt.Fprintln(os.Stderr, errorMessage)

	os.Exit(exitCode)
}

func (app *Application) GenerateDoc(dirPath string) error {
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("Error removing directory: %v\n", err)
	}
	fmt.Printf("Directory %s removed.\n", dirPath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("Error recreating directory: %v\n", err)
	}

	if err := doc.GenMarkdownTree(app.cmd.Command, dirPath); err != nil {
		return fmt.Errorf("generate doc: %w", err)
	}

	return nil
}
