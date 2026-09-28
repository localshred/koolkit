package run

import (
	"errors"
	"testing"
)

func TestRequirePresentAndMissing(t *testing.T) {
	f := &Fake{Present: map[string]bool{"gh": true}}
	if err := Require(f, "gh", "install gh"); err != nil {
		t.Fatalf("present tool errored: %v", err)
	}
	err := Require(f, "sf", "run: brew install sf")
	if err == nil {
		t.Fatal("missing tool should error")
	}
	if !contains(err.Error(), "sf") || !contains(err.Error(), "brew install sf") {
		t.Fatalf("error %q should name the tool and the hint", err)
	}
}

func TestJSONUnmarshalsStdout(t *testing.T) {
	f := &Fake{Results: map[string]Result{
		"gh api user": {Stdout: `{"login":"bj"}`},
	}}
	var got struct{ Login string }
	if err := JSON(f, &got, "gh", "api", "user"); err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if got.Login != "bj" {
		t.Fatalf("login = %q, want bj", got.Login)
	}
}

func TestFakeRunReturnsInjectedError(t *testing.T) {
	sentinel := errors.New("boom")
	f := &Fake{Errs: map[string]error{"gh pr view": sentinel}}
	if _, err := f.Run("gh", "pr", "view"); !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
