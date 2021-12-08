package initializer

import (
	"fmt"
	"github.com/ijingjingyang/go-frame/conf"
	"github.com/jjonline/go-lib-backend/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

//go:noinline
func initDB() *gorm.DB {
	c := conf.Config.Database

	if c.PoolMaxIdle <= 0 {
		c.PoolMaxIdle = 2
	}
	if c.PoolMaxOpen <= 0 {
		c.PoolMaxOpen = 5
	}
	if c.PoolMaxTime <= 0 {
		c.PoolMaxTime = 3600
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Config.Database.Prefix, // 表前缀
			SingularTable: true,                        // 关闭复数表名
		},
		PrepareStmt:            true,                    // 启用预编译
		SkipDefaultTransaction: true,                    // 对于写操作（创建、更新、删除）禁用默认事务
		Logger:                 logger.NewGorm2Logger(), // 自定义logger
	})
	if err != nil {
		panic(fmt.Sprintf("db init open err: %s", err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("db.DB() err: %s", err.Error()))
	}

	sqlDB.SetMaxIdleConns(c.PoolMaxIdle)
	sqlDB.SetMaxOpenConns(c.PoolMaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(c.PoolMaxTime) * time.Second)

	return db
}
