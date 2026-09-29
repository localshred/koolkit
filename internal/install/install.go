package install

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/localshred/koolkit/internal/registry"
	toml "github.com/pelletier/go-toml/v2"
)

type Entry struct {
	Repo        string    `toml:"repo"`
	InstalledAt time.Time `toml:"installed_at"`
}

type Manifest map[string]Entry

func manifestPath(dir string) string { return filepath.Join(dir, "manifest.toml") }

func LoadManifest(dir string) (Manifest, error) {
	data, err := os.ReadFile(manifestPath(dir))
	if os.IsNotExist(err) {
		return Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}
	var wrap struct {
		Plugins Manifest `toml:"plugins"`
	}
	if err := toml.Unmarshal(data, &wrap); err != nil {
		return nil, err
	}
	if wrap.Plugins == nil {
		wrap.Plugins = Manifest{}
	}
	return wrap.Plugins, nil
}

func (m Manifest) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := toml.Marshal(struct {
		Plugins Manifest `toml:"plugins"`
	}{Plugins: m})
	if err != nil {
		return err
	}
	return os.WriteFile(manifestPath(dir), data, 0o644)
}

// Install downloads, verifies, extracts, and atomically places kk-<name>.
func Install(dl Downloader, binDir string, p registry.Plugin, goos, goarch string) error {
	asset, err := RenderAsset(p.Asset, goos, goarch)
	if err != nil {
		return err
	}
	tgz, err := dl.Get(AssetURL(p.Host, p.Repo, asset))
	if err != nil {
		return err
	}
	sums, err := dl.Get(ChecksumsURL(p.Host, p.Repo))
	if err != nil {
		return err
	}
	if err := Verify(sums, asset, tgz); err != nil {
		return err
	}
	bin, err := extractBinary(tgz, "kk-"+p.Name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(binDir, ".kk-"+p.Name+"-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	tmp.Close()
	if err := os.Chmod(tmpName, 0o755); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, filepath.Join(binDir, "kk-"+p.Name)); err != nil {
		os.Remove(tmpName)
		return err
	}
	man, err := LoadManifest(binDir)
	if err != nil {
		return err
	}
	man[p.Name] = Entry{Repo: p.Repo, InstalledAt: time.Now().UTC()}
	return man.Save(binDir)
}

func Uninstall(binDir, name string) error {
	if err := os.Remove(filepath.Join(binDir, "kk-"+name)); err != nil && !os.IsNotExist(err) {
		return err
	}
	man, err := LoadManifest(binDir)
	if err != nil {
		return err
	}
	delete(man, name)
	return man.Save(binDir)
}

func extractBinary(tgz []byte, wantName string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(tgz))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag == tar.TypeReg && filepath.Base(h.Name) == wantName {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("binary %s not found in archive", wantName)
}
