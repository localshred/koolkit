// Package registry parses koolkit registry TOML files and resolves plugin
// names to their release source. Multiple registries merge like brew taps.
package registry

import (
	"errors"
	"io/fs"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Plugin struct {
	Name  string `toml:"name"`
	Repo  string `toml:"repo"`  // owner/repo
	Asset string `toml:"asset"` // goreleaser asset template
}

type Registry struct {
	Plugins []Plugin `toml:"plugin"`
}

func Parse(data []byte) (Registry, error) {
	var r Registry
	if err := toml.Unmarshal(data, &r); err != nil {
		return Registry{}, err
	}
	return r, nil
}

// LoadFile reads a registry file. A missing file is an empty registry (not
// an error); malformed TOML is an error.
func LoadFile(path string) (Registry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Registry{}, nil
	}
	if err != nil {
		return Registry{}, err
	}
	return Parse(data)
}

// Merge flattens registries into a name→Plugin map; earlier sources win.
func Merge(regs ...Registry) map[string]Plugin {
	out := map[string]Plugin{}
	for _, r := range regs {
		for _, p := range r.Plugins {
			if _, seen := out[p.Name]; !seen {
				out[p.Name] = p
			}
		}
	}
	return out
}

func Resolve(merged map[string]Plugin, name string) (Plugin, bool) {
	p, ok := merged[name]
	return p, ok
}
