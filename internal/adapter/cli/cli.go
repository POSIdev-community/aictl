package cli

import (
	"fmt"
	"github.com/POSIdev-community/aictl/pkg/logger"
	"io"
	"os"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/port"
)

var _ port.Cli = &Cli{}

type Cli struct {
	logger *logger.Logger
}

func NewCli(logger *logger.Logger) *Cli {
	return &Cli{logger}
}

func (cli *Cli) AskConfirmation(question string) (bool, error) {
	cli.logger.StdErr("%s [y/n]: ", question)

	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		return false, err
	}

	return strings.ToLower(answer) == "y" ||
		strings.ToLower(answer) == "yes", nil
}

func (cli *Cli) ShowProjects(projects []project.Project) {
	const format = "%-36s\t%s\n"

	cli.logger.StdErr(format, "ID", "NAME")

	for _, p := range projects {
		cli.logger.StdErr(format, p.Id, p.Name)
	}
}

func (cli *Cli) ShowProjectsQuite(projects []project.Project) {
	for _, p := range projects {
		cli.logger.StdOut(p.Id.String())
	}
}

func (cli *Cli) ShowText(text string) {
	cli.logger.StdOut(text)
}

func (cli *Cli) ShowTextF(format string, a ...any) {
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

	cli.logger.StdErr(format, "ID", "SETTINGS ID")

	for _, p := range scans {
		cli.logger.StdErr(format, p.Id, p.SettingsId)
	}
}
