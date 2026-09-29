package pluginmgr

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/localshred/koolkit/pkg/run"
)

func writeExec(t *testing.T, dir, name, body string) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

type fakeRunner struct {
	describes map[string]string
}

func (f fakeRunner) Run(name string, args ...string) (run.Result, error) {
	if len(args) == 1 && args[0] == "--kk-describe" {
		if desc, ok := f.describes[name]; ok {
			return run.Result{Stdout: desc}, nil
		}
	}
	return run.Result{}, nil
}

func (f fakeRunner) Look(name string) (string, bool) {
	_, ok := f.describes[name]
	return "/fake/path/" + name, ok
}

type slowRunner struct{}

func (slowRunner) Run(name string, args ...string) (run.Result, error) {
	time.Sleep(2 * time.Second)
	return run.Result{Stdout: "late"}, nil
}

func (slowRunner) Look(name string) (string, bool) {
	return "/fake/path/" + name, true
}

func TestDiscoverReadsDescribe(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-goto", `echo stub`)
	writeExec(t, dir, "notaplugin", `echo hi`)
	m := New(dir)
	m.Run = fakeRunner{describes: map[string]string{
		filepath.Join(dir, "kk-goto"): "navigate SCM urls",
	}}
	got, err := m.Discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "goto" || got[0].Describe != "navigate SCM urls" {
		t.Fatalf("discover = %+v", got)
	}
}

func TestDiscoverReadsDescribeViaAbsolutePath(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-real", `[ "$1" = "--kk-describe" ] && echo "real plugin desc"`)
	m := New(dir)
	got, err := m.Discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Describe != "real plugin desc" {
		t.Fatalf("describe by absolute path = %+v", got)
	}
}

func TestDescribeTimesOut(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-slow", `echo stub`)
	m := New(dir)
	m.Run = slowRunner{}
	m.Timeout = 50 * time.Millisecond
	start := time.Now()
	got, _ := m.Discover()
	elapsed := time.Since(start)
	if elapsed > 500*time.Millisecond {
		t.Fatalf("discover blocked on a hung plugin: took %v", elapsed)
	}
	if len(got) != 1 || got[0].Describe != "" {
		t.Fatalf("hung describe should yield empty desc, got %+v", got)
	}
}
