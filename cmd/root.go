/*
版权所有 © 2026 MochiZen Dan <3394549538@qq.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mocli",
	Short: "MOCLI - 命令行工具",
	Long:  `MOCLI 是一个基于 Cobra 和 Bubble Tea 的命令行工具`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
