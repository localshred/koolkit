package pluginmgr

import (
	"context"
	"os/exec"
	"sort"
	"strings"

	"github.com/localshred/koolkit/pkg/run"
)

type ToolStatus struct {
	Tool     string
	Present  bool
	NeededBy []string
}

func (m *Manager) Requires(path string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, "--kk-requires").Output()
	if err != nil {
		return nil
	}
	var tools []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if t := strings.TrimSpace(line); t != "" {
			tools = append(tools, t)
		}
	}
	return tools
}

// Doctor reports the subject tools every installed plugin requires, each
// annotated with presence on PATH and the plugins that need it.
func Doctor(m *Manager, r run.Runner) ([]ToolStatus, error) {
	plugins, err := m.Discover()
	if err != nil {
		return nil, err
	}
	needed := map[string][]string{}
	for _, p := range plugins {
		for _, tool := range m.Requires(p.Path) {
			needed[tool] = append(needed[tool], p.Name)
		}
	}
	var out []ToolStatus
	for tool, by := range needed {
		_, present := r.Look(tool)
		out = append(out, ToolStatus{Tool: tool, Present: present, NeededBy: by})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Tool < out[j].Tool })
	return out, nil
}
