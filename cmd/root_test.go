package cmd

import (
	"bytes"
	"testing"
)

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
