package pluginmgr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeExec(t *testing.T, dir, name, body string) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverReadsDescribe(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-goto", `[ "$1" = "--kk-describe" ] && echo "navigate SCM urls"`)
	writeExec(t, dir, "notaplugin", `echo hi`)
	m := New(dir)
	got, err := m.Discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "goto" || got[0].Describe != "navigate SCM urls" {
		t.Fatalf("discover = %+v", got)
	}
}

func TestDescribeTimesOut(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-slow", `[ "$1" = "--kk-describe" ] && sleep 5; echo late`)
	m := New(dir)
	m.Timeout = 200 * time.Millisecond
	start := time.Now()
	got, _ := m.Discover()
	if time.Since(start) > 2*time.Second {
		t.Fatal("discover blocked on a hung plugin")
	}
	if len(got) != 1 || got[0].Describe != "" {
		t.Fatalf("hung describe should yield empty desc, got %+v", got)
	}
}
