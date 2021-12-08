package command

import (
	"github.com/spf13/cobra"
)

// init version子命令
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "test",                              // 子命令名称
		Short: "测试命令",                              // 子命令简短说明
		Long:  "测试命令：请在下方Run方法书写测试代码，测试代码不要提交到代码库", // 子命令完整说明
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
}
