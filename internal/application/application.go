package application

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/POSIdev-community/aictl/internal/adapter/config"
	"github.com/POSIdev-community/aictl/internal/core/apperror"
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
	cmd, err := app.cmd.ExecuteContextC(ctx)
	if err == nil {
		os.Exit(ExitCodeSuccess)
	}

	reportCtx := ctx
	if cmd != nil {
		reportCtx = cmd.Context()
	}

	exitCode := reportCommandError(reportCtx, err)
	os.Exit(exitCode)
}

// reportCommandError prints a single user-facing message to stderr, mirrors it to the
// log file when configured, and emits the full Go error chain only at debug level.
// With -v/--verbose (or -V), HTTP Bad Request response bodies are also printed.
func reportCommandError(ctx context.Context, err error) int {
	log := logger.FromContext(ctx)
	exitCode, errorMessage := mapExitCode(err)

	_, _ = fmt.Fprintln(os.Stderr, errorMessage)
	log.FileError(errorMessage)

	if log.IsVerbose() {
		if detail := verboseErrorDetail(err); detail != "" {
			_, _ = fmt.Fprintln(os.Stderr, detail)
			log.FileError(detail)
		}
	}

	log.Debugf("%v", err)

	return exitCode
}

func verboseErrorDetail(err error) string {
	var badRequestErr *apperror.BadRequestError
	if errors.As(err, &badRequestErr) {
		return badRequestErr.Body()
	}

	return ""
}

func (app *Application) GenerateDoc(dirPath string) error {
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("error removing directory: %v", err)
	}
	fmt.Printf("Directory %s removed.\n", dirPath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("error recreating directory: %v", err)
	}

	// cobra GenMarkdownTree skips Deprecated commands (same as shell completion).
	// Keep them in doc/gen while they remain callable compatibility aliases.
	restore := revealDeprecatedCommandsForDocs(app.cmd.Command)
	defer restore()

	if err := doc.GenMarkdownTree(app.cmd.Command, dirPath); err != nil {
		return fmt.Errorf("generate doc: %w", err)
	}

	return nil
}
