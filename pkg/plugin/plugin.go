// Package plugin provides helpers for kk plugins: the discovery contract
// (--kk-describe / --kk-requires) that kk uses to list and doctor plugins.
package plugin

import (
	"fmt"
	"io"
)

// Meta declares a plugin's one-line description and the subject tools it
// needs at runtime.
type Meta struct {
	Describe string
	Requires []string
}

// Handle intercepts the discovery flags. Call it at the top of main with
// os.Args[1:]; if handled, print nothing further and exit with code.
func Handle(args []string, meta Meta, out io.Writer) (handled bool, code int) {
	if len(args) != 1 {
		return false, 0
	}
	switch args[0] {
	case "--kk-describe":
		fmt.Fprintln(out, meta.Describe)
		return true, 0
	case "--kk-requires":
		for _, t := range meta.Requires {
			fmt.Fprintln(out, t)
		}
		return true, 0
	}
	return false, 0
}
