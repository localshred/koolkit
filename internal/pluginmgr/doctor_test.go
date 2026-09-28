package pluginmgr

import (
	"testing"

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
