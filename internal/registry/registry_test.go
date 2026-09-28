package registry

import (
	"os"
	"path/filepath"
	"testing"
)

const sample = `
[[plugin]]
name = "goto"
repo = "localshred/kk-goto"
asset = "kk-goto_{{.Os}}_{{.Arch}}.tar.gz"

[[plugin]]
name = "prs"
repo = "localshred/kk-prs"
asset = "kk-prs_{{.Os}}_{{.Arch}}.tar.gz"
`

func TestParseAndResolve(t *testing.T) {
	reg, err := Parse([]byte(sample))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	merged := Merge(reg)
	p, ok := Resolve(merged, "goto")
	if !ok || p.Repo != "localshred/kk-goto" {
		t.Fatalf("resolve goto = %+v ok=%v", p, ok)
	}
	if _, ok := Resolve(merged, "nope"); ok {
		t.Fatal("unknown name should not resolve")
	}
}

func TestMergeFirstSourceWins(t *testing.T) {
	a := Registry{Plugins: []Plugin{{Name: "goto", Repo: "a/kk-goto"}}}
	b := Registry{Plugins: []Plugin{{Name: "goto", Repo: "b/kk-goto"}}}
	merged := Merge(a, b)
	if merged["goto"].Repo != "a/kk-goto" {
		t.Fatalf("first source should win, got %q", merged["goto"].Repo)
	}
}

func TestLoadFileMissingIsEmptyNotError(t *testing.T) {
	reg, err := LoadFile(filepath.Join(t.TempDir(), "absent.toml"))
	if err != nil || len(reg.Plugins) != 0 {
		t.Fatalf("missing file: reg=%+v err=%v", reg, err)
	}
}

func TestLoadFileMalformedErrors(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.toml")
	os.WriteFile(p, []byte("this is = = not toml ["), 0o644)
	if _, err := LoadFile(p); err == nil {
		t.Fatal("malformed TOML should error")
	}
}
