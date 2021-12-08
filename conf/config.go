package conf

import (
	"github.com/ijingjingyang/go-frame/define"
	"github.com/ijingjingyang/go-frame/utils/cfg"
)

// 项目config配置定义

// 暴露给全局使用的配置变量
var (
	Config config
	Cmd    cmdConfig
)

// config 项目配置上层结构
type config struct {
	Server   server   `json:"server"`   // 服务端口等
	Log      log      `json:"log"`      // 日志
	Database database `json:"database"` // 数据库
	Redis    redis    `json:"redis"`    // redis
}

// cmdConfig 命令行相关配置
type cmdConfig struct {
	ConfigFile  string // 命令行指定配置文件
	ConfigType  string // 命令行指定配置文件类型
	Path        string // 命令行指定日志文件，优先级高于配置文件
	Level       string // 命令行指定日志级别，优先级高于配置文件
	OnlyQueue   bool   // 命令行指定仅启动队列
	WithQueue   bool   // 命令行指定跟随启动队列
	OnlyCrontab bool   // 命令行指定仅启动定时任务
	WithCrontab bool   // 命令行指定跟随启动定时任务
}

// parseAfterLoad 配置项加载完成后的统一处理流程逻辑
func (c config) parseAfterLoad() {
	// ①、命令行的日志路径 & 日志级别高于配置文件
	if Cmd.Path != define.DefaultLogPath {
		Config.Log.Path = Cmd.Path
	}
	if Cmd.Level != define.DefaultLogLevel {
		Config.Log.Level = Cmd.Level
	}

	// ②、Url后缀斜杠处理，配置写法里带不带斜杠不影响代码使用
}

// region 初始化

// Init 初始化
func Init() {
	// must load or panic quit
	var cfgLoader cfg.IFace
	cfgLoader = cfg.Viper{}
	if err := cfgLoader.Parse(Cmd.ConfigFile, Cmd.ConfigType, &Config); err != nil {
		panic(err)
	}

	// 配置加载并解析映射成功后统一处理逻辑：譬如Url统一处理后缀斜杠
	Config.parseAfterLoad()
}

// endregion

// region 热重载

// Reload 热重载
// 当配置变更时又无需重新启动进程触发，监听调用
func Reload() {

}

// endregion
