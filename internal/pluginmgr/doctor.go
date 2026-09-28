package pluginmgr

import (
	"sort"
	"strings"
	"time"

	"github.com/localshred/koolkit/pkg/run"
)

// ToolStatus aggregates a subject tool's presence and the installed plugins
// that require it.
type ToolStatus struct {
	Tool     string
	Present  bool
	NeededBy []string
}

// Requires returns the subject tools a plugin declares via --kk-requires,
// one per line. A hung plugin is bounded by m.Timeout (returns nil).
func (m *Manager) Requires(path string) []string {
	ch := make(chan []string, 1) // buffered cap 1 — abandoned goroutine never blocks/leaks
	go func() {
		res, err := m.Run.Run(path, "--kk-requires")
		if err != nil {
			ch <- nil
			return
		}
		var tools []string
		for _, line := range strings.Split(res.Stdout, "\n") {
			if t := strings.TrimSpace(line); t != "" {
				tools = append(tools, t)
			}
		}
		ch <- tools
	}()
	select {
	case tools := <-ch:
		return tools
	case <-time.After(m.Timeout):
		return nil // hung plugin: skip its requires, Doctor keeps going
	}
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
