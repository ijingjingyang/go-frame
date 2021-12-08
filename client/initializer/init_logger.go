package initializer

import (
	"github.com/ijingjingyang/go-frame/conf"
	"github.com/jjonline/go-lib-backend/logger"
)

//go:noinline
func iniLogger() *logger.Logger {
	return logger.New(conf.Config.Log.Level, conf.Config.Log.Path)
}
