// Package xdg resolves koolkit's state and config directories per the XDG
// Base Directory spec.
package xdg

import (
	"os"
	"path/filepath"
)

func base(envVar, homeRel string) string {
	if v := os.Getenv(envVar); v != "" {
		return filepath.Join(v, "koolkit")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, filepath.FromSlash(homeRel), "koolkit")
}

// DataDir is where app-managed mutable state lives.
func DataDir() string { return base("XDG_DATA_HOME", ".local/share") }

// ConfigDir is where user settings live.
func ConfigDir() string { return base("XDG_CONFIG_HOME", ".config") }

// BinDir is the managed plugin binary directory.
func BinDir() string { return filepath.Join(DataDir(), "bin") }
