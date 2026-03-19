/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// readFileCmd represents the readFile command
var readFileCmd = &cobra.Command{
	Use:   "readFile <filePath>",
	Short: "Read a file and print its contents",
	Long: `Read the target file and print all contents to stdout.

Example:
  MOCLI readFile test.txt`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		content, err := readFileContent(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read file %q: %v\n", filePath, err)
			os.Exit(1)
		}

		fmt.Print(content)
	},
}

func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func init() {
	rootCmd.AddCommand(readFileCmd)
}
