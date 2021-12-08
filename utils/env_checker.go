package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/ijingjingyang/go-frame/conf"
)

// config.Server.env配置项可选枚举值
const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

// EnvMap env 和 运行模式对应关系
var EnvMap = map[string]string{
	EnvDev:  gin.DebugMode,
	EnvTest: gin.TestMode,
	EnvProd: gin.ReleaseMode,
}

// RunMode 获取当前配置文件中env配置项转换为gin的运行模式
func RunMode() string {
	if mode, ok := EnvMap[conf.Config.Server.Env]; ok {
		return mode
	}
	panic("config.Server.env配置项值错误")
}

// IsProd 当前是否为prod生产部署环境
func IsProd() bool {
	return conf.Config.Server.Env == EnvProd
}
