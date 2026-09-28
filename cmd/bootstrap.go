package cmd

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/localshred/koolkit/internal/config"
	"github.com/localshred/koolkit/internal/install"
	"github.com/localshred/koolkit/internal/pluginmgr"
	"github.com/localshred/koolkit/internal/registry"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/localshred/koolkit/pkg/pick"
	"github.com/spf13/cobra"
)

func newBootstrapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Fuzzy-select and install plugins from the registry",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(xdg.ConfigDir())
			if err != nil {
				return err
			}
			avail, _ := registry.Available(install.HTTP, cfg.Sources)
			installed, _ := pluginmgr.New(xdg.BinDir()).Discover()
			have := map[string]bool{}
			for _, p := range installed {
				have[p.Name] = true
			}
			var choices []registry.Plugin
			for name, p := range avail {
				if !have[name] {
					choices = append(choices, p)
				}
			}
			sort.Slice(choices, func(i, j int) bool { return choices[i].Name < choices[j].Name })
			if len(choices) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All available plugins are already installed.")
				return nil
			}
			picked, err := pick.Many(choices, func(p registry.Plugin) string {
				return fmt.Sprintf("%s  (%s)", p.Name, p.Repo)
			})
			if err != nil {
				return err
			}
			for _, p := range picked {
				if err := install.Install(install.HTTP, xdg.BinDir(), p, runtime.GOOS, runtime.GOARCH); err != nil {
					return fmt.Errorf("install %s: %w", p.Name, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "installed %s\n", p.Name)
			}
			return nil
		},
	}
}
