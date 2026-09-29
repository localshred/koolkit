package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/localshred/koolkit/internal/pluginmgr"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/spf13/cobra"
)

// Dispatch execs kk-<name> from the managed dir (or PATH), passing args and
// inheriting stdio. Returns the plugin's exit code, or non-zero with an
// error when no such plugin is installed.
func Dispatch(mgr *pluginmgr.Manager, name string, args []string, stdout, stderr io.Writer) int {
	path, ok := mgr.Find(name)
	if !ok {
		fmt.Fprintf(stderr, "kk: unknown command %q. Install it with `kk install %s`, or run `kk list`.\n", name, name)
		return 1
	}
	cmd := exec.Command(path, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, stdout, stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode()
		}
		fmt.Fprintf(stderr, "kk: failed to run %s: %v\n", name, err)
		return 1
	}
	return 0
}

// Execute routes the process args and returns an exit code. A plugin
// invocation (first arg is a bare name that is not a built-in command) is
// dispatched raw so the plugin's own flags never reach kk's flag parser;
// everything else runs through cobra.
func Execute(version string) int {
	root := NewRoot(version)
	if name, rest, ok := pluginInvocation(root, os.Args[1:]); ok {
		mgr := pluginmgr.New(xdg.BinDir())
		return Dispatch(mgr, name, rest, os.Stdout, os.Stderr)
	}
	if err := root.Execute(); err != nil {
		return 1
	}
	return 0
}

// NewRoot builds the kk root command. Built-in subcommands are attached by
// their own files (install, list, registry, bootstrap, doctor); plugin
// dispatch is handled ahead of cobra in Execute.
func NewRoot(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "kk",
		Short:         "koolkit — a self-contained dev-utility plugin suite",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.AddCommand(newBootstrapCmd(), newDoctorCmd(), newInstallCmd(), newUninstallCmd(), newListCmd(), newRegistryCmd())
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	return root
}

// pluginInvocation reports whether args target a plugin rather than a
// built-in command, splitting the plugin name from the args to forward. The
// first arg must be a bare name (no leading dash) that is not a registered
// subcommand or alias; a leading flag or a known command falls through to
// cobra.
func pluginInvocation(root *cobra.Command, args []string) (name string, rest []string, ok bool) {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return "", nil, false
	}
	for _, c := range root.Commands() {
		if c.Name() == args[0] || c.HasAlias(args[0]) {
			return "", nil, false
		}
	}
	return args[0], args[1:], true
}
