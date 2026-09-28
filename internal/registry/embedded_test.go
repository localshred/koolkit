// koolkit/internal/registry/embedded_test.go
package registry

import (
	"os"
	"testing"
)

func TestShippedRegistryParses(t *testing.T) {
	data, err := os.ReadFile("../../registry.toml")
	if err != nil {
		t.Fatal(err)
	}
	reg, err := Parse(data)
	if err != nil {
		t.Fatalf("shipped registry.toml is malformed: %v", err)
	}
	merged := Merge(reg)
	for _, name := range []string{"goto", "prs"} {
		if _, ok := merged[name]; !ok {
			t.Errorf("registry missing %q", name)
		}
	}
}
