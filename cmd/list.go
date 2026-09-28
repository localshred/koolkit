package cmd

import (
	"github.com/localshred/koolkit/internal/config"
	"github.com/localshred/koolkit/internal/install"
	"github.com/localshred/koolkit/internal/pluginmgr"
	"github.com/localshred/koolkit/internal/registry"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/localshred/koolkit/pkg/ui"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed and available plugins",
		RunE: func(cmd *cobra.Command, _ []string) error {
			installed, _ := pluginmgr.New(xdg.BinDir()).Discover()
			have := map[string]bool{}
			rows := [][]string{{ui.Dim("NAME"), ui.Dim("STATUS"), ui.Dim("DETAIL")}}
			for _, p := range installed {
				have[p.Name] = true
				rows = append(rows, []string{ui.Cyan(p.Name), ui.Green("installed"), ui.Dim(p.Describe)})
			}
			cfg, _ := config.Load(xdg.ConfigDir())
			avail, _ := registry.Available(install.HTTP, cfg.Sources)
			for name, p := range avail {
				if !have[name] {
					rows = append(rows, []string{ui.Cyan(name), ui.Dim("available"), ui.Dim(p.Repo)})
				}
			}
			ui.PrintTable(cmd.OutOrStdout(), rows)
			return nil
		},
	}
}
