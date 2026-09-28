// Package run executes subprocesses and detects subject tools. The Runner
// interface lets callers substitute a Fake in tests so external tools are
// not required in CI.
package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Result struct {
	Stdout string
	Stderr string
	Code   int
}

type Runner interface {
	Run(name string, args ...string) (Result, error)
	Look(name string) (string, bool)
}

type execRunner struct{}

func (execRunner) Run(name string, args ...string) (Result, error) {
	cmd := exec.Command(name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	res := Result{Stdout: out.String(), Stderr: errb.String()}
	if ee, ok := err.(*exec.ExitError); ok {
		res.Code = ee.ExitCode()
	}
	return res, err
}

func (execRunner) Look(name string) (string, bool) {
	p, err := exec.LookPath(name)
	return p, err == nil
}

// Default is the real subprocess runner.
var Default Runner = execRunner{}

// Require returns nil if name is on PATH, else an actionable error carrying
// the caller-supplied install hint.
func Require(r Runner, name, hint string) error {
	if _, ok := r.Look(name); ok {
		return nil
	}
	return fmt.Errorf("%s not found on PATH — %s", name, hint)
}

// JSON runs name+args and unmarshals stdout into v.
func JSON(r Runner, v any, name string, args ...string) error {
	res, err := r.Run(name, args...)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(res.Stdout), v)
}
