package cmd

import (
	"bytes"
	"testing"
)

func TestPluginInvocationForwardsFlags(t *testing.T) {
	root := NewRoot("0")
	name, rest, ok := pluginInvocation(root, []string{"goto", "-u", "foo", "--dry"})
	if !ok || name != "goto" {
		t.Fatalf("goto should dispatch as a plugin, got name=%q ok=%v", name, ok)
	}
	if want := []string{"-u", "foo", "--dry"}; !equalArgs(rest, want) {
		t.Fatalf("forwarded args = %v, want %v", rest, want)
	}
}

func TestPluginInvocationSkipsBuiltins(t *testing.T) {
	root := NewRoot("0")
	for _, arg := range []string{"install", "list", "help", "--version", "-h"} {
		if _, _, ok := pluginInvocation(root, []string{arg, "x"}); ok {
			t.Fatalf("%q must route to cobra, not plugin dispatch", arg)
		}
	}
	if _, _, ok := pluginInvocation(root, nil); ok {
		t.Fatal("no args must not dispatch a plugin")
	}
}

func TestRootVersion(t *testing.T) {
	root := NewRoot("1.2.3")
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--version"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if got := out.String(); got == "" || !bytes.Contains(out.Bytes(), []byte("1.2.3")) {
		t.Fatalf("version output = %q, want it to contain 1.2.3", got)
	}
}

func equalArgs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
