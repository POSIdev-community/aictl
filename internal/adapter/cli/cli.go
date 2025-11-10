package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/POSIdev-community/aictl/pkg/logger"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
)

type Cli struct {
	logger *logger.Logger
}

func NewCli(logger *logger.Logger) *Cli {
	return &Cli{logger}
}

func (cli *Cli) AskConfirmation(question string) (bool, error) {
	cli.logger.StdOutF("%s [y/n]: ", question)

	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		return false, err
	}

	return strings.ToLower(answer) == "y" ||
		strings.ToLower(answer) == "yes", nil
}

func (cli *Cli) ShowProjects(projects []project.Project) {
	const format = "%-36s\t%s"

	cli.logger.StdOutF(format, "ID", "NAME")

	for _, p := range projects {
		cli.logger.StdOutF(format, p.Id, p.Name)
	}
}

func (cli *Cli) ShowProjectsQuite(projects []project.Project) {
	for _, p := range projects {
		cli.logger.StdOut(p.Id.String())
	}
}

func (cli *Cli) ShowText(text string) {
	cli.logger.StdErr(text)
}

func (cli *Cli) ShowTextF(format string, a ...any) {
	cli.logger.StdErrF(format, a...)
}

func (cli *Cli) ReturnText(text string) {
	cli.logger.StdOut(text)
}

func (cli *Cli) ReturnTextF(format string, a ...any) {
	cli.logger.StdOutF(format, a...)
}

// ShowReader copy provided reader to stdout.
func (cli *Cli) ShowReader(r io.Reader) error {
	if _, err := io.Copy(os.Stdout, r); err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	return nil
}

func (cli *Cli) ShowScans(scans []scan.Scan) {
	const format = "%-36s\t%-36s\n"

	cli.logger.StdErrF(format, "ID", "SETTINGS ID")

	for _, p := range scans {
		cli.logger.StdErrF(format, p.Id, p.SettingsId)
	}
}
