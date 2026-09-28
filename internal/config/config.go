// Package config persists koolkit user settings (registry sources) at
// ~/.config/koolkit/config.toml.
package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

const DefaultPublicSource = "https://raw.githubusercontent.com/localshred/koolkit/main/registry.toml"

// Config holds koolkit user configuration.
type Config struct {
	Sources []string `toml:"sources"`
}

// AddSource appends src if not already present, returns whether the set changed.
func (c *Config) AddSource(src string) bool {
	for _, s := range c.Sources {
		if s == src {
			return false
		}
	}
	c.Sources = append(c.Sources, src)
	return true
}

// Load reads config from dir/config.toml. Missing file returns Config seeded with DefaultPublicSource.
func Load(dir string) (Config, error) {
	data, err := os.ReadFile(path(dir))
	if os.IsNotExist(err) {
		return Config{Sources: []string{DefaultPublicSource}}, nil
	}
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := toml.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	if len(c.Sources) == 0 {
		c.Sources = []string{DefaultPublicSource}
	}
	return c, nil
}

// RemoveSource removes src if present, returns whether the set changed.
func (c *Config) RemoveSource(src string) bool {
	for i, s := range c.Sources {
		if s == src {
			c.Sources = append(c.Sources[:i], c.Sources[i+1:]...)
			return true
		}
	}
	return false
}

// Save writes config to dir/config.toml.
func (c Config) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path(dir), data, 0o644)
}

func path(dir string) string { return filepath.Join(dir, "config.toml") }
