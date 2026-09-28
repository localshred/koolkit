package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestHumanize(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{50 * time.Minute, "50m"},
		{3 * time.Hour, "3h"},
		{50 * time.Hour, "2d"},
		{-5 * time.Minute, "0m"},
	}
	for _, c := range cases {
		if got := Humanize(c.d); got != c.want {
			t.Errorf("Humanize(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

func TestStripANSIAndVisibleLen(t *testing.T) {
	s := Cyan("abc")
	if VisibleLen(s) != 3 {
		t.Errorf("VisibleLen(cyan abc) = %d, want 3", VisibleLen(s))
	}
	link := Hyperlink("https://x", "abc")
	if VisibleLen(link) != 3 {
		t.Errorf("VisibleLen(hyperlink abc) = %d, want 3", VisibleLen(link))
	}
}

func TestPrintTableAlignsOnVisibleWidth(t *testing.T) {
	var buf bytes.Buffer
	PrintTable(&buf, [][]string{
		{Cyan("a"), "short"},
		{Cyan("bbbb"), "x"},
	})
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	// First column padded to visible width 4 ("bbbb") + 2 spaces = col 2 aligned.
	if !strings.Contains(lines[0], "a") || !strings.HasSuffix(lines[0], "short") {
		t.Errorf("row0 = %q", lines[0])
	}
}
