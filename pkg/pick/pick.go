// Package pick wraps a fuzzy finder for single and multi selection,
// replacing the external fzf dependency.
package pick

import (
	"errors"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// Finder abstracts the underlying picker so tests can inject selections.
type Finder func(n int, label func(i int) string, multi bool) ([]int, error)

var find Finder = func(n int, label func(i int) string, multi bool) ([]int, error) {
	if multi {
		return fuzzyfinder.FindMulti(make([]struct{}, n),
			func(i int) string { return label(i) })
	}
	i, err := fuzzyfinder.Find(make([]struct{}, n),
		func(i int) string { return label(i) })
	if err != nil {
		return nil, err
	}
	return []int{i}, nil
}

// ErrCancelled is returned when the user aborts the picker.
var ErrCancelled = fuzzyfinder.ErrAbort

// One prompts for a single selection. ok is false when cancelled.
func One[T any](items []T, label func(T) string) (T, bool, error) {
	var zero T
	idx, err := find(len(items), func(i int) string { return label(items[i]) }, false)
	if errors.Is(err, ErrCancelled) {
		return zero, false, nil
	}
	if err != nil {
		return zero, false, err
	}
	return items[idx[0]], true, nil
}

// Many prompts for a multi selection.
func Many[T any](items []T, label func(T) string) ([]T, error) {
	idxs, err := find(len(items), func(i int) string { return label(items[i]) }, true)
	if errors.Is(err, ErrCancelled) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(idxs))
	for _, i := range idxs {
		out = append(out, items[i])
	}
	return out, nil
}
