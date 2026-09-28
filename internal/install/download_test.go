package install

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

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
