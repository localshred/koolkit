package registry

import (
	"strings"
)

// Downloader fetches registry content from URLs. Defined locally to avoid an
// import cycle with internal/install. Consumers pass install.HTTP, which
// structurally satisfies this interface.
type Downloader interface {
	Get(url string) ([]byte, error)
}

// Available loads and merges every enabled source. A URL source is fetched
// via dl; anything else is read as a local file. Bad sources are collected
// and skipped, never fatal.
func Available(dl Downloader, sources []string) (map[string]Plugin, []error) {
	var regs []Registry
	var errs []error
	for _, src := range sources {
		var (
			reg Registry
			err error
		)
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			var data []byte
			if data, err = dl.Get(src); err == nil {
				reg, err = Parse(data)
			}
		} else {
			reg, err = LoadFile(src)
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		regs = append(regs, reg)
	}
	return Merge(regs...), errs
}
