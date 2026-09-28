package plugin

import (
	"bytes"
	"strings"
	"testing"
)

func TestHandleDescribe(t *testing.T) {
	var out bytes.Buffer
	handled, code := Handle([]string{"--kk-describe"}, Meta{Describe: "does the thing"}, &out)
	if !handled || code != 0 {
		t.Fatalf("handled=%v code=%d", handled, code)
	}
	if strings.TrimSpace(out.String()) != "does the thing" {
		t.Fatalf("out = %q", out.String())
	}
}

func TestHandleRequiresOnePerLine(t *testing.T) {
	var out bytes.Buffer
	Handle([]string{"--kk-requires"}, Meta{Requires: []string{"gh", "op"}}, &out)
	if lines := strings.Fields(out.String()); len(lines) != 2 || lines[0] != "gh" || lines[1] != "op" {
		t.Fatalf("out = %q", out.String())
	}
}

func TestHandlePassesThroughNormalArgs(t *testing.T) {
	var out bytes.Buffer
	if handled, _ := Handle([]string{"-u", "foo"}, Meta{Describe: "x"}, &out); handled {
		t.Fatal("normal args should not be handled")
	}
}
