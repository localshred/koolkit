// Package ui provides terminal-UI helpers: ANSI colors, OSC 8 hyperlinks,
// width-aware table rendering, and compact duration formatting. Ported
// from the mx babashka ui.clj.
package ui

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
	ansiDim    = "\033[2m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiReset  = "\033[0m"
	ansiYellow = "\033[33m"
)

var (
	reSGR  = regexp.MustCompile(`\033\[[0-9;]*m`)
	reOSC8 = regexp.MustCompile("\033\\]8;;[^\033]*\033\\\\")
)

func colorize(color, s string) string { return color + s + ansiReset }

func Blue(s string) string   { return colorize(ansiBlue, s) }
func Cyan(s string) string   { return colorize(ansiCyan, s) }
func Dim(s string) string    { return colorize(ansiDim, s) }
func Green(s string) string  { return colorize(ansiGreen, s) }
func Red(s string) string    { return colorize(ansiRed, s) }
func Yellow(s string) string { return colorize(ansiYellow, s) }

// Hyperlink wraps text in an OSC 8 terminal hyperlink (iTerm2, WezTerm, …).
func Hyperlink(url, text string) string {
	return "\033]8;;" + url + "\033\\" + text + "\033]8;;\033\\"
}

// Humanize renders a duration compactly as days, hours, or minutes.
func Humanize(d time.Duration) string {
	if days := int(d.Hours()) / 24; days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	if hours := int(d.Hours()); hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	mins := int(d.Minutes())
	if mins < 0 {
		mins = 0
	}
	return fmt.Sprintf("%dm", mins)
}

// AgeColor returns an ANSI color by staleness: red past 14 days, yellow past 7.
func AgeColor(d time.Duration) string {
	days := int(d.Hours()) / 24
	switch {
	case days >= 14:
		return ansiRed
	case days >= 7:
		return ansiYellow
	default:
		return ansiDim
	}
}

// AgeCell is a colored, humanized duration cell.
func AgeCell(d time.Duration) string { return AgeColor(d) + Humanize(d) + ansiReset }

// StripANSI removes OSC 8 hyperlink and SGR color escapes.
func StripANSI(s string) string { return reSGR.ReplaceAllString(reOSC8.ReplaceAllString(s, ""), "") }

// VisibleLen is the on-screen rune count of a colored/linked cell.
func VisibleLen(s string) int { return len([]rune(StripANSI(s))) }

func pad(cell string, width int) string {
	n := width - VisibleLen(cell)
	if n < 0 {
		n = 0
	}
	return cell + strings.Repeat(" ", n)
}

// PrintTable prints rows as columns aligned on visible width; the last
// column of each row is left free (not padded).
func PrintTable(w io.Writer, rows [][]string) {
	cols := 0
	for _, r := range rows {
		if len(r) > cols {
			cols = len(r)
		}
	}
	widths := make([]int, cols)
	for _, r := range rows {
		for i, c := range r {
			if l := VisibleLen(c); l > widths[i] {
				widths[i] = l
			}
		}
	}
	for _, r := range rows {
		cells := make([]string, len(r))
		for i, c := range r {
			if i == len(r)-1 {
				cells[i] = c
			} else {
				cells[i] = pad(c, widths[i])
			}
		}
		fmt.Fprintln(w, strings.Join(cells, "  "))
	}
}
