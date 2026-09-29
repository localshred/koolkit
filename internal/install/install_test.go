package install

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/localshred/koolkit/internal/registry"
)

func makeTarGz(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg})
	tw.Write(content)
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

type mapDownloader map[string][]byte

func (m mapDownloader) Get(url string) ([]byte, error) {
	if b, ok := m[url]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("404 %s", url)
}

func TestInstallPlacesVerifiedBinary(t *testing.T) {
	binDir := t.TempDir()
	p := registry.Plugin{Name: "goto", Repo: "localshred/kk-goto", Asset: "kk-goto_{{.Os}}_{{.Arch}}.tar.gz"}
	tgz := makeTarGz(t, "kk-goto", []byte("BINARY"))
	sum := sha256.Sum256(tgz)
	sums := fmt.Sprintf("%s  kk-goto_Darwin_arm64.tar.gz\n", hex.EncodeToString(sum[:]))
	dl := mapDownloader{
		AssetURL(p.Host, p.Repo, "kk-goto_Darwin_arm64.tar.gz"): tgz,
		ChecksumsURL(p.Host, p.Repo):                            []byte(sums),
	}
	if err := Install(dl, binDir, p, "darwin", "arm64"); err != nil {
		t.Fatalf("install: %v", err)
	}
	placed := filepath.Join(binDir, "kk-goto")
	b, err := os.ReadFile(placed)
	if err != nil || string(b) != "BINARY" {
		t.Fatalf("placed binary = %q err=%v", b, err)
	}
	fi, _ := os.Stat(placed)
	if fi.Mode().Perm()&0o100 == 0 {
		t.Fatal("placed binary not executable")
	}
	man, _ := LoadManifest(binDir)
	if man["goto"].Repo != p.Repo {
		t.Fatalf("manifest missing goto: %+v", man)
	}
}

func TestInstallChecksumMismatchLeavesNoBinary(t *testing.T) {
	binDir := t.TempDir()
	p := registry.Plugin{Name: "goto", Repo: "localshred/kk-goto", Asset: "kk-goto_{{.Os}}_{{.Arch}}.tar.gz"}
	tgz := makeTarGz(t, "kk-goto", []byte("BINARY"))
	dl := mapDownloader{
		AssetURL(p.Host, p.Repo, "kk-goto_Darwin_arm64.tar.gz"): tgz,
		ChecksumsURL(p.Host, p.Repo):                            []byte("deadbeef  kk-goto_Darwin_arm64.tar.gz\n"),
	}
	if err := Install(dl, binDir, p, "darwin", "arm64"); err == nil {
		t.Fatal("mismatch should error")
	}
	if _, err := os.Stat(filepath.Join(binDir, "kk-goto")); !os.IsNotExist(err) {
		t.Fatal("no binary should be left after a failed verify")
	}
}

func TestInstall404Aborts(t *testing.T) {
	binDir := t.TempDir()
	p := registry.Plugin{Name: "goto", Repo: "localshred/kk-goto", Asset: "kk-goto_{{.Os}}_{{.Arch}}.tar.gz"}
	if err := Install(mapDownloader{}, binDir, p, "darwin", "arm64"); err == nil {
		t.Fatal("missing asset should error")
	}
	if _, err := os.Stat(filepath.Join(binDir, "kk-goto")); !os.IsNotExist(err) {
		t.Fatal("no binary should be left after 404")
	}
}
