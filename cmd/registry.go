package cmd

import (
	"fmt"

	"github.com/localshred/koolkit/internal/config"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/spf13/cobra"
)

func newRegistryCmd() *cobra.Command {
	c := &cobra.Command{Use: "registry", Short: "Manage registry sources"}
	c.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List registry sources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(xdg.ConfigDir())
			if err != nil {
				return err
			}
			for _, s := range cfg.Sources {
				fmt.Fprintln(cmd.OutOrStdout(), s)
			}
			return nil
		},
	})
	c.AddCommand(&cobra.Command{
		Use:   "add <source>",
		Short: "Add a registry source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(xdg.ConfigDir())
			if err != nil {
				return err
			}
			if cfg.AddSource(args[0]) {
				return cfg.Save(xdg.ConfigDir())
			}
			fmt.Fprintln(cmd.OutOrStdout(), "already present")
			return nil
		},
	})
	c.AddCommand(&cobra.Command{
		Use:   "remove <source>",
		Short: "Remove a registry source",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(xdg.ConfigDir())
			if err != nil {
				return err
			}
			if cfg.RemoveSource(args[0]) {
				return cfg.Save(xdg.ConfigDir())
			}
			fmt.Fprintln(cmd.OutOrStdout(), "not found")
			return nil
		},
	})
	return c
}
