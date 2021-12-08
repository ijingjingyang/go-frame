package client

// 定义全局句柄、需要初始化的变量，然后在initializer子包内单个文件进行初始化

import (
	"github.com/go-redis/redis/v8"
	"github.com/jjonline/go-lib-backend/guzzle"
	"github.com/jjonline/go-lib-backend/logger"
	"github.com/jjonline/go-lib-backend/memory"
	"github.com/jjonline/go-lib-backend/queue"
	"gorm.io/gorm"
)

// 全局使用的句柄变量--client
//  ① 项目启动就要初始化的全局变量句柄
//  ② 按client包规则书写
var (
	Logger      *logger.Logger // 基于zap的日志记录器
	Redis       *redis.Client  // redis客户端
	DB          *gorm.DB       // 数据库客户端
	MemoryCache *memory.Cache  // 本地内存缓存客户端
	Guzzle      *guzzle.Client // http客户端简单封装
	Queue       *queue.Queue   // 队列
)
