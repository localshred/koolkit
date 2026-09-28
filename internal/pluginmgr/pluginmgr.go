// Package pluginmgr discovers installed kk-* plugin binaries and reads
// their --kk-describe one-liners with a timeout guard.
package pluginmgr

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const prefix = "kk-"

type Plugin struct {
	Name     string
	Path     string
	Describe string
}

type Manager struct {
	Dir     string
	Timeout time.Duration
}

func New(dir string) *Manager { return &Manager{Dir: dir, Timeout: 2 * time.Second} }

// Find returns the absolute path of kk-<name>, preferring the managed dir
// then falling back to PATH (development fallback).
func (m *Manager) Find(name string) (string, bool) {
	p := filepath.Join(m.Dir, prefix+name)
	if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
		return p, true
	}
	if p, err := exec.LookPath(prefix + name); err == nil {
		return p, true
	}
	return "", false
}

// Discover lists kk-* executables in the managed dir, sorted by name, each
// with its --kk-describe line (empty on error/timeout).
func (m *Manager) Discover() ([]Plugin, error) {
	entries, err := os.ReadDir(m.Dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Plugin
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), prefix) {
			continue
		}
		path := filepath.Join(m.Dir, e.Name())
		out = append(out, Plugin{
			Name:     strings.TrimPrefix(e.Name(), prefix),
			Path:     path,
			Describe: m.describe(path),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Manager) describe(path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, "--kk-describe").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
