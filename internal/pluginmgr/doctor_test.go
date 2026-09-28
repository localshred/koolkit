package pluginmgr

import (
	"testing"
	"time"

	"github.com/localshred/koolkit/pkg/run"
)

func TestDoctorAggregatesRequires(t *testing.T) {
	dir := t.TempDir()
	writeExec(t, dir, "kk-prs", `[ "$1" = "--kk-requires" ] && printf 'gh\n'`)
	writeExec(t, dir, "kk-widget", `[ "$1" = "--kk-requires" ] && printf 'gh\nk9s\n'`)
	m := New(dir)
	statuses, err := Doctor(m, &run.Fake{Present: map[string]bool{"gh": true}})
	if err != nil {
		t.Fatal(err)
	}
	byTool := map[string]ToolStatus{}
	for _, s := range statuses {
		byTool[s.Tool] = s
	}
	if !byTool["gh"].Present || len(byTool["gh"].NeededBy) != 2 {
		t.Fatalf("gh status = %+v", byTool["gh"])
	}
	if byTool["k9s"].Present || len(byTool["k9s"].NeededBy) != 1 {
		t.Fatalf("k9s status = %+v", byTool["k9s"])
	}
}

type blockingRunner struct{}

func (blockingRunner) Run(name string, args ...string) (run.Result, error) {
	time.Sleep(200 * time.Millisecond)
	return run.Result{Stdout: "late"}, nil
}

func (blockingRunner) Look(name string) (string, bool) {
	return "/fake/path/" + name, true
}

func TestRequiresTimesOutHungPlugin(t *testing.T) {
	m := New(t.TempDir())
	m.Run = blockingRunner{}
	m.Timeout = 50 * time.Millisecond
	start := time.Now()
	got := m.Requires("/fake/hung-plugin")
	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Fatalf("Requires blocked on hung plugin: took %v", elapsed)
	}
	if got != nil {
		t.Fatalf("hung plugin should return nil, got %v", got)
	}
}
