package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

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

// NewRoot builds the kk root command. Subcommands are attached by their
// own files (install, list, registry, bootstrap, doctor); unknown args
// fall through to plugin dispatch (see Task A9).
func NewRoot(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "kk",
		Short:         "koolkit — a self-contained dev-utility plugin suite",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.Args = cobra.ArbitraryArgs
	root.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		mgr := pluginmgr.New(xdg.BinDir())
		if code := Dispatch(mgr, args[0], args[1:], cmd.OutOrStdout(), cmd.ErrOrStderr()); code != 0 {
			os.Exit(code)
		}
		return nil
	}
	root.AddCommand(newBootstrapCmd(), newDoctorCmd(), newInstallCmd(), newUninstallCmd(), newListCmd(), newRegistryCmd())
	return root
}

// Execute runs the root command and returns a process exit code.
func Execute(version string) int {
	if err := NewRoot(version).Execute(); err != nil {
		return 1
	}
	return 0
}
