/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"MOCLI/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var readFileCmd = &cobra.Command{
	Use:   "rf <filePath>",
	Short: "Read a file and print its contents",
	Long: `Read the target file and print all contents to stdout.

Example:
  mocli rf test.txt`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		content, err := internal.ReadFileContent(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", filePath, err)
		}

		fmt.Fprint(cmd.OutOrStdout(), content)
		return nil
	},
}
