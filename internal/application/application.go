package application

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/POSIdev-community/aictl/internal/adapter/config"
	"github.com/POSIdev-community/aictl/internal/di"
	"github.com/POSIdev-community/aictl/internal/presenter"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type Application struct {
	cmd *presenter.CmdRoot
}

func NewApplication() (*Application, error) {
	cfgAdapter := config.NewContextAdapter()
	cfg := cfgAdapter.GetContextFromAictlFolder()

	cmd, err := di.InitializeCmd(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize commands: %w", err)
	}
	cmd.DisableAutoGenTag = true

	return &Application{cmd}, nil
}

func (app *Application) Run(ctx context.Context) {
	err := app.cmd.ExecuteContext(ctx)
	if err == nil {
		os.Exit(ExitCodeSuccess)
	}
	log := logger.FromContext(ctx)
	log.StdErrf(err.Error())

	exitCode, errorMessage := mapExitCode(err)

	_, _ = fmt.Fprintln(os.Stderr, errorMessage)

	os.Exit(exitCode)
}

func (app *Application) GenerateDoc(dirPath string) error {
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("error removing directory: %v", err)
	}
	fmt.Printf("Directory %s removed.\n", dirPath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("error recreating directory: %v", err)
	}

	if err := doc.GenMarkdownTree(app.cmd.Command, dirPath); err != nil {
		return fmt.Errorf("generate doc: %w", err)
	}

	return nil
}
