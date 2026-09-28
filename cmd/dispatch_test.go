package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/localshred/koolkit/internal/pluginmgr"
)

func TestDispatchExecsPlugin(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "kk-echo"),
		[]byte("#!/bin/sh\necho \"$@\"\n"), 0o755)
	var out bytes.Buffer
	code := Dispatch(pluginmgr.New(dir), "echo", []string{"hello", "world"}, &out, &out)
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if got := out.String(); got != "hello world\n" {
		t.Fatalf("plugin output = %q", got)
	}
}

func TestDispatchUnknownPluginErrors(t *testing.T) {
	var out, errb bytes.Buffer
	code := Dispatch(pluginmgr.New(t.TempDir()), "bogus", nil, &out, &errb)
	if code == 0 {
		t.Fatal("unknown plugin should return non-zero")
	}
	if !bytes.Contains(errb.Bytes(), []byte("bogus")) {
		t.Fatalf("stderr %q should name the unknown command", errb.String())
	}
}
