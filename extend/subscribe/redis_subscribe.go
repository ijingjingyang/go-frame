package subscribe

import (
	"context"
	"github.com/ijingjingyang/go-frame/client"
	"github.com/ijingjingyang/go-frame/utils/memory"
	"strings"
)

// RedisSubscribe redis 发布-订阅功能实现的多pod间触发执行一些同步逻辑
//  - 极端场景下redis pub-sub可靠性并不能百分百保障订阅者取得每一份消息
//  - <ugly-极端场景发生时业务已不可用订阅者是否还能收到消息还重要？>
//  - go-redis 本身已内部实现断线自动重连 https://redis.uptrace.dev/#pubsub
type RedisSubscribe struct{}

// New 实例化redis订阅同步实例
func New() *RedisSubscribe {
	return &RedisSubscribe{}
}

// Start 启动redis发布-订阅监听<内部协程方式>
func (rs RedisSubscribe) Start() {
	go func() {
		// 监听订阅取得 chan
		sub := client.Redis.Subscribe(context.Background(), "--your-redis-pub-sub-key--")
		ch := sub.Channel()

		// 轮询 chan 获得监听得到的消息
		for msg := range ch {
			// 通过消息本身的内容调度不同逻辑（发布的消息作为命令标识）
			switch msg.Payload {
			case EventCleanAppMemCache:
				rs.cleanAppMemoryCache()
			}

			client.Logger.Debug("subscribe message occur")
		}
	}()
}

// cleanAppMemoryCache sub-pub模式各pod同步清理本地mem内存
func (rs RedisSubscribe) cleanAppMemoryCache() {
	// 清理App相关本地内存缓存
	items := client.MemoryCache.Items()
	for key := range items {
		if strings.HasPrefix(key, "-- your memory cache key --") {
			memory.Del(key)
			client.Logger.Info(key) // log info for subscribe event
		}
	}
}
