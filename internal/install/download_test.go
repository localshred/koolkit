package install

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestAssetURLDefaultHost(t *testing.T) {
	got := AssetURL("", "localshred/kk-goto", "kk-goto_Darwin_arm64.tar.gz")
	want := "https://github.com/localshred/kk-goto/releases/latest/download/kk-goto_Darwin_arm64.tar.gz"
	if got != want {
		t.Errorf("AssetURL default host = %q", got)
	}
	if sums := ChecksumsURL("", "localshred/kk-goto"); !strings.HasPrefix(sums, "https://github.com/") {
		t.Errorf("ChecksumsURL default host = %q", sums)
	}
}

func TestAssetURLExplicitHost(t *testing.T) {
	got := AssetURL("ghe.example.com", "alice/kk-tool", "kk-tool_Linux_amd64.tar.gz")
	want := "https://ghe.example.com/alice/kk-tool/releases/latest/download/kk-tool_Linux_amd64.tar.gz"
	if got != want {
		t.Errorf("AssetURL explicit host = %q", got)
	}
}

func TestRenderAsset(t *testing.T) {
	got, err := RenderAsset("kk-goto_{{.Os}}_{{.Arch}}.tar.gz", "darwin", "arm64")
	if err != nil {
		t.Fatal(err)
	}
	if got != "kk-goto_Darwin_arm64.tar.gz" {
		t.Fatalf("RenderAsset = %q", got)
	}
}

func TestVerifyMatches(t *testing.T) {
	data := []byte("payload")
	sum := sha256.Sum256(data)
	line := fmt.Sprintf("%s  kk-goto_Darwin_arm64.tar.gz\n", hex.EncodeToString(sum[:]))
	if err := Verify([]byte(line), "kk-goto_Darwin_arm64.tar.gz", data); err != nil {
		t.Fatalf("Verify should pass: %v", err)
	}
}

func TestVerifyMismatchErrors(t *testing.T) {
	line := "deadbeef  kk-goto_Darwin_arm64.tar.gz\n"
	if err := Verify([]byte(line), "kk-goto_Darwin_arm64.tar.gz", []byte("payload")); err == nil {
		t.Fatal("checksum mismatch should error")
	}
}

func TestVerifyMissingEntryErrors(t *testing.T) {
	if err := Verify([]byte("x  other.tar.gz\n"), "kk-goto_Darwin_arm64.tar.gz", nil); err == nil {
		t.Fatal("missing checksum entry should error")
	}
}
