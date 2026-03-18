/*
版权所有 © 2026 MochiZen Dan <3394549538@qq.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 表示基础命令（当没有任何子命令时调用）
var rootCmd = &cobra.Command{
	Use:   "MOCLI",
	Short: "design for ",
	Long:  `MOCLI is Commendline which can used for `,

	// 如果你的根命令需要执行逻辑，可以在这里定义运行函数
	Run: func(cmd *cobra.Command, args []string) {
		toggle, _ := cmd.Flags().GetBool("toggle")
		if toggle {
			fmt.Println("toggle")
		}
		fmt.Printf("Hello World\n")
	},
}

// Execute 用于将所有子命令添加到根命令，并正确设置参数
// 该函数会在 main.main() 中被调用
// 对 rootCmd 来说只需要调用一次
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 在这里定义参数（flags）和配置项

	// Cobra 支持持久化参数（Persistent Flags），
	// 如果在这里定义，则对整个应用全局生效（所有子命令都可以使用）

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径（默认是 $HOME/.MOCLI.yaml）")

	// Cobra 也支持本地参数（Local Flags），
	// 这些参数只在当前命令被直接调用时生效

	rootCmd.Flags().BoolP("toggle", "t", false, "toggle 参数的帮助信息")
}
