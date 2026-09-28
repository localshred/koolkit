package config

import "testing"

func TestLoadMissingReturnsDefault(t *testing.T) {
	c, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Sources) != 1 || c.Sources[0] != DefaultPublicSource {
		t.Fatalf("default sources = %v", c.Sources)
	}
}

func TestAddRemoveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	c, _ := Load(dir)
	if !c.AddSource("/tmp/internal.toml") {
		t.Fatal("AddSource should report a change")
	}
	if c.AddSource("/tmp/internal.toml") {
		t.Fatal("duplicate AddSource should report no change")
	}
	if err := c.Save(dir); err != nil {
		t.Fatal(err)
	}
	got, _ := Load(dir)
	if len(got.Sources) != 2 {
		t.Fatalf("reloaded sources = %v", got.Sources)
	}
	if !got.RemoveSource("/tmp/internal.toml") {
		t.Fatal("RemoveSource should report a change")
	}
}
