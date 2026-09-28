package run

import "strings"

// Fake is a scripted Runner for tests. Results/Errs are keyed by the
// space-joined command line ("gh api user").
type Fake struct {
	Results map[string]Result
	Errs    map[string]error
	Present map[string]bool
	Calls   [][]string
}

func fakeKey(name string, args []string) string {
	return strings.TrimSpace(name + " " + strings.Join(args, " "))
}

func (f *Fake) Run(name string, args ...string) (Result, error) {
	f.Calls = append(f.Calls, append([]string{name}, args...))
	k := fakeKey(name, args)
	if err := f.Errs[k]; err != nil {
		return Result{}, err
	}
	return f.Results[k], nil
}

func (f *Fake) Look(name string) (string, bool) {
	if f.Present[name] {
		return "/fake/" + name, true
	}
	return "", false
}
