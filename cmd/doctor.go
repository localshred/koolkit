package cmd

import (
	"fmt"

	"github.com/localshred/koolkit/internal/pluginmgr"
	"github.com/localshred/koolkit/internal/xdg"
	"github.com/localshred/koolkit/pkg/run"
	"github.com/localshred/koolkit/pkg/ui"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Report subject-tool availability for installed plugins",
		RunE: func(cmd *cobra.Command, _ []string) error {
			statuses, err := pluginmgr.Doctor(pluginmgr.New(xdg.BinDir()), run.Default)
			if err != nil {
				return err
			}
			if len(statuses) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No installed plugins require external tools.")
				return nil
			}
			rows := [][]string{{ui.Dim("TOOL"), ui.Dim("STATUS"), ui.Dim("NEEDED BY")}}
			for _, s := range statuses {
				status := ui.Green("✓ found")
				if !s.Present {
					status = ui.Red("✗ missing")
				}
				rows = append(rows, []string{ui.Cyan(s.Tool), status, ui.Dim(joinComma(s.NeededBy))})
			}
			ui.PrintTable(cmd.OutOrStdout(), rows)
			return nil
		},
	}
}

func joinComma(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
