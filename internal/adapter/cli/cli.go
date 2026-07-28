package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/POSIdev-community/aictl/internal/core/domain/branch"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/queue"
	"github.com/POSIdev-community/aictl/internal/core/domain/report"
	"github.com/POSIdev-community/aictl/internal/core/domain/scan"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanagent"
	"github.com/POSIdev-community/aictl/internal/core/domain/settings"
	"github.com/POSIdev-community/aictl/internal/core/domain/statistic"
	"github.com/POSIdev-community/aictl/pkg/logger"
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

func (cli *Adapter) ShowBranches(ctx context.Context, branches []branch.Branch) {
	log := logger.FromContext(ctx)
	const format = "%-36s\t%s"

	log.StdOutf(format, "ID", "NAME")

	for _, b := range branches {
		log.StdOutf(format, b.Id, b.Name)
	}
}

func (cli *Adapter) ShowBranchesQuite(ctx context.Context, branches []branch.Branch) {
	log := logger.FromContext(ctx)

	for _, b := range branches {
		log.StdOut(b.Id.String())
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

	const format = "%-36s\t%-24s\t%s"

	log.StdOutf(format, "ID", "DATE", "LABEL")

	for _, s := range scans {
		date := ""
		if !s.ScanDate.IsZero() {
			date = s.ScanDate.Format("2006-01-02 15:04:05")
		}

		log.StdOutf(format, s.Id, date, s.ScanLabel)
	}
}

func (cli *Adapter) ShowScansQuite(ctx context.Context, scans []scan.Scan) {
	log := logger.FromContext(ctx)

	for _, s := range scans {
		log.StdOut(s.Id.String())
	}
}

func (cli *Adapter) ShowScanStatistic(ctx context.Context, statistic *statistic.Statistic) {
	log := logger.FromContext(ctx)

	log.StdOutf("Total: %d", statistic.Total)
	log.StdOutf("High: %d", statistic.High)
	log.StdOutf("Medium: %d", statistic.Medium)
	log.StdOutf("Low: %d", statistic.Low)
	log.StdOutf("Potential: %d", statistic.Potential)
}

func (cli *Adapter) ShowScanAgents(ctx context.Context, agents []scanagent.ScanAgent) {
	log := logger.FromContext(ctx)
	const format = "%-36s\t%-24s\t%-12s\t%-12s\t%s"

	log.StdOutf(format, "ID", "NAME", "STATUS", "VERSION", "OS")

	for _, a := range agents {
		log.StdOutf(format, a.Id, a.Name, a.Status, a.Version, a.OperatingSystem)
	}
}

func (cli *Adapter) ShowScanAgentsQuite(ctx context.Context, agents []scanagent.ScanAgent) {
	log := logger.FromContext(ctx)

	for _, a := range agents {
		log.StdOut(a.Id.String())
	}
}

func (cli *Adapter) ShowQueue(ctx context.Context, entries []queue.Entry) {
	log := logger.FromContext(ctx)
	const format = "%-36s\t%-36s\t%-36s\t%s"

	log.StdOutf(format, "PROJECT ID", "BRANCH ID", "SCAN ID", "STAGE")

	for _, e := range entries {
		log.StdOutf(format, e.ProjectId, e.BranchId, e.ScanId, e.Stage)
	}
}

func (cli *Adapter) ShowReportTemplates(ctx context.Context, templates []report.Template) {
	log := logger.FromContext(ctx)
	const format = "%-36s\t%s"

	log.StdOutf(format, "ID", "NAME")

	for _, t := range templates {
		log.StdOutf(format, t.Id, t.Name)
	}
}

func (cli *Adapter) ShowReportTemplatesQuite(ctx context.Context, templates []report.Template) {
	log := logger.FromContext(ctx)

	for _, t := range templates {
		log.StdOut(t.Id.String())
	}
}

func (cli *Adapter) ShowProjectSettings(ctx context.Context, view settings.ProjectSettingsView) {
	log := logger.FromContext(ctx)

	log.StdOutf("Priority: %s", view.Priority)
	log.StdOutf("Preferred agents only: %t", view.PreferredAgentsOnly)

	if len(view.PreferredAgents) == 0 {
		log.StdOut("Preferred agents: none")

		return
	}

	log.StdOut("Preferred agents:")
	for _, id := range view.PreferredAgents {
		log.StdOutf("  %s", id)
	}
}
