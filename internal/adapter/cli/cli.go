package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/POSIdev-community/aictl/pkg/logger"

	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (cli *Adapter) AskConfirmation(ctx context.Context, question string) (bool, error) {
	log := logger.FromContext(ctx)
	log.StdOutf("%s [y/n]: ", question)

	var answer string
	_, err := fmt.Scan(&answer)
	if err != nil {
		return false, err
	}

	return strings.ToLower(answer) == "y" ||
		strings.ToLower(answer) == "yes", nil
}

func (cli *Adapter) ShowProjects(ctx context.Context, projects []project.Project) {
	log := logger.FromContext(ctx)
	const format = "%-36s\t%s"

	log.StdOutf(format, "ID", "NAME")

	for _, p := range projects {
		log.StdOutf(format, p.Id, p.Name)
	}
}

func (cli *Adapter) ShowProjectsQuite(ctx context.Context, projects []project.Project) {
	log := logger.FromContext(ctx)

	for _, p := range projects {
		log.StdOut(p.Id.String())
	}
}

func (cli *Adapter) ShowText(ctx context.Context, text string) {
	log := logger.FromContext(ctx)

	log.StdErr(text)
}

func (cli *Adapter) ShowTextf(ctx context.Context, format string, a ...any) {
	log := logger.FromContext(ctx)

	log.StdErrf(format, a...)
}

func (cli *Adapter) ReturnText(ctx context.Context, text string) {
	log := logger.FromContext(ctx)

	log.StdOut(text)
}

func (cli *Adapter) ReturnTextf(ctx context.Context, format string, a ...any) {
	log := logger.FromContext(ctx)

	log.StdOutf(format, a...)
}

// ShowReader copy provided reader to stdout.
func (cli *Adapter) ShowReader(r io.Reader) error {
	if _, err := io.Copy(os.Stdout, r); err != nil {
		return fmt.Errorf("write to stdout: %w", err)
	}

	return nil
}

func (cli *Adapter) ShowScans(ctx context.Context, scans []scan.Scan) {
	log := logger.FromContext(ctx)

	const format = "%-36s\t%-36s\n"

	log.StdErrf(format, "ID", "SETTINGS ID")

	for _, p := range scans {
		log.StdErrf(format, p.Id, p.SettingsId)
	}
}
