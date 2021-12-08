package command

import (
	"fmt"
	"github.com/ijingjingyang/go-frame/app/console"
	"github.com/ijingjingyang/go-frame/client/initializer"
	"github.com/ijingjingyang/go-frame/conf"
	"github.com/ijingjingyang/go-frame/define"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd 基于cobra的命令行根节点定义
var (
	RootCmd = &cobra.Command{
		Use:   "sufficient [sub]",
		Short: "sufficient服务管理",
		Long: `sufficient服务管理，源码模式下参数样例如下:
----------------------------------------------------------------------
go run main.go               仅启动基础api服务
go run main.go --withQueue   启动基础api服务的同时启动队列消费者服务
go run main.go --withCrontab 启动基础api服务的同时启动定时服务
go run main.go --onlyQueue   不启动基础api服务仅启动队列消费者服务
go run main.go --onlyCrontab 不启动基础api服务仅启动队定定时任务
go run main.go --log=stderr  命令行参数指定优先级高于配置文件的日志路径
go run main.go --level=info  命令行参数指定优先级高于配置文件的日志级别
go run main.go customSubCmd  自定义子命令
----------------------------------------------------------------------
编译后的二进制文件请更换上述"go run main.go"为对应可执行二进制文件名即可`,
		Run: func(cmd *cobra.Command, args []string) {
			console.BootStrap()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// step1、init config
			conf.Init()

			// step2、init global client handle
			initializer.Init()

			return nil
		},
	}
)

func init() {
	// 全局配置
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.ConfigFile, "config", "conf.toml", "指定本地配置文件")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.ConfigType, "configType", "toml", "指定配置文件类型")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.Path, "log", define.DefaultLogPath, "指定优先级高于配置文件的日志存储位置：stderr|stdout|目录路径")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.Level, "level", define.DefaultLogLevel, "指定优先级高于配置文件的日志级别：debug|info|warn|error|panic|fatal")

	// 命令配置
	RootCmd.Flags().BoolVar(&conf.Cmd.WithCrontab, "withCrontab", false, "跟随启动定时任务")
	RootCmd.Flags().BoolVar(&conf.Cmd.OnlyCrontab, "onlyCrontab", false, "仅启动定时任务")
	RootCmd.Flags().BoolVar(&conf.Cmd.WithQueue, "withQueue", false, "跟随启动队列消费者服务")
	RootCmd.Flags().BoolVar(&conf.Cmd.OnlyQueue, "onlyQueue", false, "仅启动队列消费者服务")
}

// Start 启动应用
func Start() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
