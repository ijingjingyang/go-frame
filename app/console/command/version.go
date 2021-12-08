package command

import (
	"fmt"
	"github.com/ijingjingyang/go-frame/conf"
	"github.com/spf13/cobra"
)

// init version子命令
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "version", // 子命令名称
		Short: "显示当前版本号", // 子命令简短说明
		Long:  "显示当前版本号", // 子命令完整说明
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(conf.Config.Server.Version)
		},
	})
}
