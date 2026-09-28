package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/localshred/koolkit/internal/config"
	"github.com/localshred/koolkit/internal/install"
	"github.com/localshred/koolkit/internal/registry"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <name|github.com/owner/kk-name>...",
		Short: "Install one or more plugins",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(xdg.ConfigDir())
			if err != nil {
				return err
			}
			avail, srcErrs := registry.Available(install.HTTP, cfg.Sources)
			for _, e := range srcErrs {
				fmt.Fprintf(cmd.ErrOrStderr(), "kk: skipping registry source: %v\n", e)
			}
			for _, name := range args {
				p, err := resolvePlugin(avail, name)
				if err != nil {
					return err
				}
				if err := install.Install(install.HTTP, xdg.BinDir(), p, runtime.GOOS, runtime.GOARCH); err != nil {
					return fmt.Errorf("install %s: %w", p.Name, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "installed %s\n", p.Name)
			}
			return nil
		},
	}
}

func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <name>...",
		Short: "Remove installed plugins",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				if err := install.Uninstall(xdg.BinDir(), name); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "removed %s\n", name)
			}
			return nil
		},
	}
}

// resolvePlugin returns a registry entry for name, or synthesizes one from a
// direct github path (github.com/owner/kk-foo).
func resolvePlugin(avail map[string]registry.Plugin, name string) (registry.Plugin, error) {
	if strings.HasPrefix(name, "github.com/") {
		repo := strings.TrimPrefix(name, "github.com/")
		base := repo[strings.LastIndex(repo, "/")+1:]
		short := strings.TrimPrefix(base, "kk-")
		return registry.Plugin{Name: short, Repo: repo, Asset: base + "_{{.Os}}_{{.Arch}}.tar.gz"}, nil
	}
	if p, ok := registry.Resolve(avail, name); ok {
		return p, nil
	}
	return registry.Plugin{}, fmt.Errorf("no plugin %q in any registry; try a direct github.com/owner/kk-%s path", name, name)
}
