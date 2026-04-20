package cmd

import (
	"MOCLI/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read [file] or read -i [startPath]",
	Short: "读取文件内容",
	Args: func(cmd *cobra.Command, args []string) error {
		interactive, _ := cmd.Flags().GetBool("interactive")
		if interactive {
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive, _ := cmd.Flags().GetBool("interactive")

		var path string
		if interactive {
			startPath := "."
			if len(args) > 0 {
				startPath = args[0]
			}
			selected, err := app.SelectFile(startPath)
			if err != nil {
				return err
			}
			if selected == "" {
				return nil
			}
			path = selected
		} else {
			path = args[0]
		}

		content, err := app.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), content)
		return nil
	},
}

func init() {
	readCmd.Flags().BoolP("interactive", "i", false, "交互式文件选择")
	rootCmd.AddCommand(readCmd)
}
