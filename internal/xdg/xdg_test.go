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
