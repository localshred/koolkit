package xdg

import (
	"path/filepath"
	"testing"
)

func TestDataDirHonorsXDG(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/tmp/xdgdata")
	if got, want := DataDir(), "/tmp/xdgdata/koolkit"; got != want {
		t.Fatalf("DataDir = %q, want %q", got, want)
	}
	if got, want := BinDir(), "/tmp/xdgdata/koolkit/bin"; got != want {
		t.Fatalf("BinDir = %q, want %q", got, want)
	}
}

func TestDataDirFallsBackToHome(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("HOME", "/home/bj")
	if got, want := DataDir(), filepath.FromSlash("/home/bj/.local/share/koolkit"); got != want {
		t.Fatalf("DataDir = %q, want %q", got, want)
	}
}

func TestConfigDirHonorsXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdgconfig")
	if got, want := ConfigDir(), filepath.Join("/tmp/xdgconfig", "koolkit"); got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
}

func TestConfigDirFallsBackToHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/home/bj")
	if got, want := ConfigDir(), filepath.Join("/home/bj", filepath.FromSlash(".config"), "koolkit"); got != want {
		t.Fatalf("ConfigDir = %q, want %q", got, want)
	}
}
