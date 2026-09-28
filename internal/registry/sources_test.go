package registry

import (
	"os"
	"path/filepath"
	"testing"
)

type fakeDL map[string][]byte

func (f fakeDL) Get(url string) ([]byte, error) {
	if b, ok := f[url]; ok {
		return b, nil
	}
	return nil, os.ErrNotExist
}

func TestAvailableMergesGoodSkipsBad(t *testing.T) {
	dir := t.TempDir()
	local := filepath.Join(dir, "public.toml")
	os.WriteFile(local, []byte("[[plugin]]\nname=\"goto\"\nrepo=\"localshred/kk-goto\"\n"), 0o644)
	dl := fakeDL{"https://x/reg.toml": []byte("[[plugin]]\nname=\"prs\"\nrepo=\"localshred/kk-prs\"\n")}

	merged, errs := Available(dl, []string{local, "https://x/reg.toml", "https://x/missing.toml"})
	if _, ok := merged["goto"]; !ok {
		t.Error("goto from local source missing")
	}
	if _, ok := merged["prs"]; !ok {
		t.Error("prs from url source missing")
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 source error (the missing url), got %d", len(errs))
	}
}
