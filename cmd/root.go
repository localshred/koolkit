package cmd

import "github.com/spf13/cobra"

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
	return root
}

// Execute runs the root command and returns a process exit code.
func Execute(version string) int {
	if err := NewRoot(version).Execute(); err != nil {
		return 1
	}
	return 0
}
