package cmd

import (
	"MOCLI/tui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open Bubble Tea TUI",
	Run: func(cmd *cobra.Command, args []string) {
		tui.Run()
	},
}
