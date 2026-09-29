// Package install resolves, downloads, verifies, and places plugin
// binaries in the managed dir.
package install

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type Downloader interface {
	Get(url string) ([]byte, error)
}

type httpDownloader struct{}

func (httpDownloader) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// HTTP is the real downloader.
var HTTP Downloader = httpDownloader{}

func AssetURL(repo, asset string) string {
	return fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", repo, asset)
}

func ChecksumsURL(repo string) string {
	return fmt.Sprintf("https://github.com/%s/releases/latest/download/checksums.txt", repo)
}

// RenderAsset expands a goreleaser asset template. Os is title-cased to
// match goreleaser defaults (Darwin/Linux); Arch stays lowercase.
func RenderAsset(tmpl, goos, goarch string) (string, error) {
	t, err := template.New("asset").Parse(tmpl)
	if err != nil {
		return "", err
	}
	data := struct{ Os, Arch string }{Os: titleFirst(goos), Arch: goarch}
	var sb strings.Builder
	if err := t.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// Verify checks data's sha256 against the checksums.txt entry for filename.
func Verify(sums []byte, filename string, data []byte) error {
	want := ""
	for _, line := range strings.Split(string(sums), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == filename {
			want = fields[0]
			break
		}
	}
	if want == "" {
		return fmt.Errorf("no checksum entry for %s", filename)
	}
	got := sha256.Sum256(data)
	if hex.EncodeToString(got[:]) != want {
		return fmt.Errorf("checksum mismatch for %s", filename)
	}
	return nil
}

// titleFirst upper-cases the first rune of s, matching goreleaser's Os
// casing (darwin -> Darwin, linux -> Linux) without the deprecated
// strings.Title.
func titleFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
