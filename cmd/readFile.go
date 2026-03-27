/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"MOCLI/internal/fileio"
	"MOCLI/tui"
	"fmt"

	"github.com/spf13/cobra"
)

var readFileCmd = &cobra.Command{
	Use:   "rf <filePath> or rf -g [startPath]",
	Short: "Read a file and print its contents",
	Long: `Read the target file and print all contents to stdout.

Example:
  mocli rf test.txt
  mocli rf -g
  mocli rf -g .`,

	Args: func(cmd *cobra.Command, args []string) error {
		graphics, _ := cmd.Flags().GetBool("graphics")
		if graphics {
			if len(args) > 1 {
				return fmt.Errorf("-g 模式最多接收一个路径参数")
			}
			return nil
		}

		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		graphics, _ := cmd.Flags().GetBool("graphics")
		if graphics {
			startPath := "."
			if len(args) == 1 {
				startPath = args[0]
			}

			selectedPath, err := tui.RunFileExplorer(startPath)
			if err != nil {
				return fmt.Errorf("failed to open file explorer: %w", err)
			}

			if selectedPath == "" {
				return nil
			}

			content, err := fileio.ReadFileContent(selectedPath)
			if err != nil {
				return fmt.Errorf("failed to read file %q: %w", selectedPath, err)
			}

			fmt.Fprint(cmd.OutOrStdout(), content)
			return nil
		}

		filePath := args[0]
		content, err := fileio.ReadFileContent(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", filePath, err)
		}

		fmt.Fprint(cmd.OutOrStdout(), content)
		return nil
	},
}

func init() {
	readFileCmd.Flags().BoolP("graphics", "g", false, "Open table-based file explorer")
}
